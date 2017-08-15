package main

import (
	"os"

	"gitlab.alexhamlin.co/go/slackbridge/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
