package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/isfonzar/slack-grand-race/pkg/config"
	"github.com/isfonzar/slack-grand-race/pkg/logs"
	"github.com/isfonzar/slack-grand-race/pkg/message"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/kelseyhightower/envconfig"
	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

const (
	databaseConnectionAttempts         = 5
	databaseWaitingTimeBetweenAttempts = 3 * time.Second
)

func main() {
	fmt.Println("Starting Slack Grand Race")

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

	// Handlers
	// @todo one big handler to handle everything?
	msgHandler := message.NewHandler(logger)

	// Slack
	// @todo improve and abstract this whole slack thingy
	api := slack.New(conf.SlackToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				msg := message.NewMessageFromEvent(ev)

				err := msgHandler.HandleMessage(msg)
				if err != nil {
					logger.Error(err)
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
