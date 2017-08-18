package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"gitlab.alexhamlin.co/go/slackbridge/slackio"

	"github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:     "exec [flags] -- program [args]",
	Example: "exec -c C12345678 -- cat  # echo server",
	Short:   "connect a program's standard streams to a single Slack channel",
	Long: `exec runs a provided executable program and connects its standard
input, output, and error streams to a single Slack channel. In this mode,
text from the main body of the channel (i.e. excluding threads) is received
by the executable on stdin. Text emitted on stdout and stderr is batched over
a short time interval and sent as a single Slack message.`,

	Args: cobra.MinimumNArgs(1),
	Run:  runExecCmd,
}

func init() {
	RootCmd.AddCommand(execCmd)
	execCmd.Flags().StringP("channel", "c", "", "ID of the channel to connect to (required)")
	execCmd.MarkFlagRequired("channel")
}

func runExecCmd(cmd *cobra.Command, args []string) {
	apiToken := os.Getenv("SLACK_TOKEN")
	if apiToken == "" {
		fmt.Fprintln(os.Stderr, "Error: SLACK_TOKEN environment variable not set")
		fmt.Fprintln(os.Stderr, RootCmd.UsageString())
		os.Exit(1)
	}

	slackChannel, err := cmd.Flags().GetString("channel")
	if err != nil || slackChannel == "" {
		fmt.Fprintln(os.Stderr, "Error: requires --channel flag")
		fmt.Fprintln(os.Stderr, cmd.UsageString())
		os.Exit(1)
	}

	client := &slackio.Client{APIToken: apiToken}
	reader := &slackio.Reader{Client: client, SlackChannelID: slackChannel}
	writer := &slackio.Writer{Client: client, SlackChannelID: slackChannel}

	child := exec.Command(args[0], args[1:]...)
	child.Stdin = reader
	child.Stdout = writer
	child.Stderr = writer

	ensure := func(fn func() error) {
		if err := fn(); err != nil {
			panic(err)
		}
	}

	ensure(child.Start)

	// We want exec to terminate when its child program exits. However, child
	// maintains an internal goroutine that copies from Stdin, and the call to
	// Wait blocks on its completion. So before we call Wait we have to shut down
	// slackIO to trigger an EOF on Read. But of course we need child to exit
	// first! Blocking on SIGCHLD is semi-hackish but does the job.

	sigchld := make(chan os.Signal)
	signal.Notify(sigchld, syscall.SIGCHLD)
	<-sigchld

	ensure(writer.Close)
	ensure(reader.Close)
	ensure(client.Close)
	ensure(child.Wait)
}
