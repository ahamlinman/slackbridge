package main

import (
	"log"
	"os"

	"github.com/nlopes/slack"
)

func main() {
	var slackToken string
	var ok bool
	if slackToken, ok = os.LookupEnv("SLACK_API_TOKEN"); !ok {
		panic("SLACK_API_TOKEN not defined")
	}

	// Totally copying the example for now

	api := slack.New(slackToken)

	logger := log.New(os.Stdout, "slackbridge: ", 0)
	slack.SetLogger(logger)
	api.SetDebug(true)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for msg := range rtm.IncomingEvents {
		logger.Println(msg)
	}
}
