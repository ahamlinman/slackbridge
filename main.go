package main

import (
	"encoding/json"
	"io/ioutil"
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
	configJSON, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var config slackbridgeConfig
	err = json.Unmarshal(configJSON, &config)
	if err != nil {
		panic(err)
	}

	return config
}
