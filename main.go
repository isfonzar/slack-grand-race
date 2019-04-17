package main

import (
	"fmt"

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
}
