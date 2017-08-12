package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"gitlab.alexhamlin.co/go/slackbridge/slackio"
)

type slackbridgeConfig struct {
	APIToken string
	Channel  string
	Exec     []string
}

func main() {
	config := getConfig("./config.json")

	slackIO := slackio.New(config.APIToken, config.Channel)
	defer slackIO.Close()

	cmd := exec.Command(config.Exec[0], config.Exec[1:]...)

	cmd.Stdin = slackIO
	cmd.Stdout = slackIO
	cmd.Stderr = slackIO

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	sigchld := make(chan os.Signal)
	signal.Notify(sigchld, syscall.SIGCHLD)
	<-sigchld
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
