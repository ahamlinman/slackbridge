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

	cmd := exec.Command(config.Exec[0], config.Exec[1:]...)
	cmd.Stdin = slackIO
	cmd.Stdout = slackIO
	cmd.Stderr = slackIO

	require(cmd.Start)

	// We want slackbridge to terminate when its child program exits. However,
	// cmd maintains an internal goroutine that copies from Stdin, and the call
	// to Wait blocks on its completion. So before we call Wait we have to shut
	// down slackIO to trigger an EOF on Read. But of course we need the child to
	// exit first! Blocking on SIGCHLD is semi-hackish but does the job.

	sigchld := make(chan os.Signal)
	signal.Notify(sigchld, syscall.SIGCHLD)
	<-sigchld

	require(slackIO.Close)
	require(cmd.Wait)
}

func require(f func() error) {
	if err := f(); err != nil {
		panic(err)
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
