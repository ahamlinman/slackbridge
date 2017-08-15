package main

import "gitlab.alexhamlin.co/go/slackbridge/cmd"

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		panic(err)
	}
}
