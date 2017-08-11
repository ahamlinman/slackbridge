package main

import (
	"encoding/json"
	"os"
)

type slackbridgeConfig struct {
	APIToken string
	Channel  string
	Exec     []string
}

func main() {
	config := getConfig("./config.json")

	sls := newSlackLineStreamer(config.APIToken, config.Channel)
	defer sls.Close()

	els := newExecLineStreamer(config.Exec)
	defer els.Close()

	for {
		select {
		case line := <-sls.ReceiveChan():
			els.Send(line)

		case line := <-els.ReceiveChan():
			sls.Send(line)
		}
	}
}

func getConfig(filename string) slackbridgeConfig {
	configFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()

	var config slackbridgeConfig
	err = json.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}
