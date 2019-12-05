package main

import (
	"fmt"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/nlopes/slack"

	"github.com/isfonzar/slack-grand-race/internal/config"
	"github.com/isfonzar/slack-grand-race/pkg/logs"
)

func main() {
	fmt.Println("Starting Slack Grand Race")

	// Configuration
	conf, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}

	// Logger
	log, err := logs.New(conf.Debug)
	if err != nil {
		panic(err)
	}
	log.Debug("Logger loaded")

	fields := []interface{}{"config", conf}
	log.Debugw("Configs loaded", fields...)

	// Running migrations
	m, err := migrate.New(
		"file://db/migrations",
		fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", conf.DB.User, conf.DB.Password, conf.DB.Host, conf.DB.DatabaseName))
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

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
				fields := []interface{}{"message", ev}
				log.Debugw("Message received", fields...)
			case *slack.RTMError:
				fields := []interface{}{"error", ev.Error()}
				log.Errorw("Error", fields...)

			case *slack.InvalidAuthEvent:
				log.Fatal("Invalid auth")

			default:
				fields := []interface{}{"event", ev}
				log.Debugw("Event received", fields...)
			}
		}
	}
}
