package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"gitlab.alexhamlin.co/go/slackbridge/slackio"
)

var (
	slackChannel string
	showUsage    bool
)

func init() {
	flag.StringVar(&slackChannel, "channel", "", "Slack channel ID (required)")
	flag.BoolVar(&showUsage, "help", false, "Show usage information")

	flag.Usage = printUsage
}

func main() {
	flag.Parse()

	if showUsage {
		flag.Usage()
		return
	}

	apiToken := os.Getenv("SLACK_TOKEN")
	if apiToken == "" {
		fmt.Fprintln(os.Stderr, "required env var not provided: SLACK_TOKEN")
		flag.Usage()
		os.Exit(222)
	}

	if slackChannel == "" {
		fmt.Fprintln(os.Stderr, "required flag not provided: -channel")
		flag.Usage()
		os.Exit(222)
	}

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "required arguments not provided: [program]")
		flag.Usage()
		os.Exit(222)
	}

	slackIO := &slackio.Client{
		APIToken:       apiToken,
		SlackChannelID: slackChannel,
	}

	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)
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

func printUsage() {
	fmt.Fprintln(os.Stderr, `slackbridge connects executable programs to Slack channels

Usage

  slackbridge [flags] [program] [arguments...]

Flags
`)

	flag.PrintDefaults()
}
