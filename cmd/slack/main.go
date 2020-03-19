package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/isfonzar/slack-grand-race/internal/repository/postgres"
	"github.com/isfonzar/slack-grand-race/pkg/config"
	"github.com/isfonzar/slack-grand-race/pkg/domain"
	"github.com/isfonzar/slack-grand-race/pkg/handlers/coins"
	"github.com/isfonzar/slack-grand-race/pkg/handlers/message"
	"github.com/isfonzar/slack-grand-race/pkg/handlers/user"
	slackHandler "github.com/isfonzar/slack-grand-race/pkg/infrastructure/slack"
	"github.com/isfonzar/slack-grand-race/pkg/logs"
	"github.com/kelseyhightower/envconfig"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

const (
	databaseConnectionAttempts         = 5
	databaseWaitingTimeBetweenAttempts = 3 * time.Second
)

func main() {
	// Configuration
	conf, err := config.LoadEnv(envconfig.Process)
	if err != nil {
		panic(err)
	}

	// Logger
	logger, err := logs.New(conf.Debug)
	if err != nil {
		panic(err)
	}
	logger.Debug("Logger loaded")
	logger.Debugw("Configs loaded",
		"config", conf,
	)

	// Connecting to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.DB.Host, conf.DB.User, conf.DB.Password, conf.DB.DatabaseName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatalw("could not connect to postgres database",
			"error", err,
			"host", conf.DB.Host,
			"database", conf.DB.DatabaseName,
			"user", conf.DB.User,
		)
	}
	defer func() {
		if err = db.Close(); err != nil {
			logger.Fatal(err)
		}
	}()

	// Pinging database
	checkDatabaseConnection(db, conf, logger, databaseConnectionAttempts)

	// Running migrations
	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.DatabaseName))
	if err != nil {
		logger.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Fatal(err)
	}

	// Slack
	api := slack.New(conf.SlackToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// Storages
	userStorage := postgres.NewUsersStorage(db)

	// Handlers
	slackCommunicator := slackHandler.NewHandler(rtm)

	coinHandler := coins.NewHandler(userStorage, slackCommunicator)
	msgHandler := message.NewHandler(coinHandler, logger)
	userHandler := user.NewHandler(rtm, userStorage)

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				// If user does not exist (like spoiler messages)
				if ev.Msg.User == "" {
					continue
				}

				// Get message
				m := domain.NewMessageFromSlack(ev)

				// Get user
				u, err := userHandler.GetUser(ev.Msg.User)
				if err != nil {
					logger.Fatalw("could not get user",
						"error", err,
						"user_id", ev.Msg.User,
						"user", u,
					)
				}

				if err := msgHandler.Process(m, u); err != nil {
					logger.Fatalw("could not process message",
						"error", err,
						"message", m,
						"user", u,
					)
				}
			case *slack.RTMError:
				fields := []interface{}{"error", ev.Error()}
				logger.Errorw("Error", fields...)

			case *slack.InvalidAuthEvent:
				logger.Fatal("Invalid auth")

			default:
				fields := []interface{}{"event", ev}
				logger.Debugw("Event received", fields...)
			}
		}
	}
}

func checkDatabaseConnection(db *sql.DB, c *config.Config, logger *zap.SugaredLogger, attempts int) {
	if attempts == 0 {
		logger.Fatalw("database connection could not be established",
			"host", c.DB.Host,
			"database", c.DB.DatabaseName,
			"user", c.DB.User,
		)

		return
	}

	if err := db.Ping(); err != nil {
		logger.Infow("database connection could not be established, waiting 5 seconds to try again",
			"error", err,
			"host", c.DB.Host,
			"database", c.DB.DatabaseName,
			"user", c.DB.User,
		)

		time.Sleep(databaseWaitingTimeBetweenAttempts)
		checkDatabaseConnection(db, c, logger, attempts-1)
		return
	}
}
