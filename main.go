package main

import (
	"fmt"
	"github.com/nlopes/slack"

	"github.com/isfonzar/slack-grand-race/internal/config"
	"github.com/isfonzar/slack-grand-race/pkg/logs"
)

func main() {
	fmt.Println("Starting Slack Grand Race")

	conf, err := config.LoadEnv()
	if err != nil {
		panic(err)
	}

	log, err := logs.New(conf.Debug)
	if err != nil {
		panic(err)
	}
	log.Debug("Logger loaded")

	fields := []interface{}{"config", conf}
	log.Debugw("Configs loaded", fields...)

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
