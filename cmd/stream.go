package cmd

import (
	"fmt"
	"io"
	"os"

	"gitlab.alexhamlin.co/go/slackbridge/slackio"

	"github.com/spf13/cobra"
)

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "stream the output of a single Slack channel to stdout",
	Long: `stream connects to Slack and continuously streams the main body of a
single channel (i.e. excluding threads) to standard output.`,
	Run: runStreamCmd,
}

func init() {
	RootCmd.AddCommand(streamCmd)
	streamCmd.Flags().StringP("channel", "c", "", "ID of the channel to connect to (required)")
	streamCmd.MarkFlagRequired("channel")
}

func runStreamCmd(cmd *cobra.Command, args []string) {
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
	defer client.Close()

	slackIO := &slackio.Reader{Client: client, SlackChannelID: slackChannel}
	defer slackIO.Close()

	if _, err := io.Copy(os.Stdout, slackIO); err != nil {
		panic(err)
	}
}
