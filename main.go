package main

import (
	"fmt"
	"github.com/isfonzar/instafarm/pkg/logs"
)

func main() {
	fmt.Println("Starting Slack Grand Race")

	log, err := logs.New()
	if err != nil {
		panic(err)
	}

	log.Infow("Logger loaded")
}
