package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type slackbridgeConfig struct {
	APIToken string
	Channel  string
}

func main() {
	config := getConfig("./config.json")

	sls := newSlackLineStreamer(config.APIToken, config.Channel)
	defer sls.Close()

	for line := range sls.ReceiveChan() {
		fmt.Printf("got line: %s\n", line)
		sls.Send("ed is the\n*standard*\ntext editor")
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
