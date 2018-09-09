package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.alexhamlin.co/slackio"

	"go.alexhamlin.co/slackbridge/internal/childproc"
)

var execCmd = &cobra.Command{
	Use:     "exec [flags] -- program [args]",
	Example: "exec -c C12345678 -- cat  # echo server",
	Short:   "Connect a program's standard streams to a single Slack channel",
	Long: `Exec runs a provided executable program and connects its standard
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
	if err != nil {
		// If the flag isn't provided, Cobra already prints a nice message for us
		panic(err)
	}

	client := slackio.NewClient(apiToken)
	reader := slackio.NewReader(client, slackChannel)
	writer := slackio.NewWriter(client, slackChannel, nil)

	child, err := childproc.Spawn(args, reader, writer)
	if err != nil {
		panic(err)
	}

	// Note that Wait will close reader and writer for us after the child process
	// terminates
	if err := child.Wait(); err != nil {
		panic(err)
	}

	if err := client.Close(); err != nil {
		panic(err)
	}
}
