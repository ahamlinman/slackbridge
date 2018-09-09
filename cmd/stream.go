package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"go.alexhamlin.co/slackio"
)

var streamCmd = &cobra.Command{
	Use:   "stream",
	Short: "Stream one or more Slack channels to stdout",
	Long: `Stream connects to Slack and continuously streams the main body of one
or more channels (i.e. excluding threads) to standard output. By default, the
text of all of the user's channels will be streamed together with no
identification of any message's originating channel. If desired, output can be
filtered to a single channel.`,
	Run: runStreamCmd,
}

func init() {
	RootCmd.AddCommand(streamCmd)
	streamCmd.Flags().StringP("channel", "c", "", "only output messages from the provided channel ID")
}

func runStreamCmd(cmd *cobra.Command, args []string) {
	apiToken := os.Getenv("SLACK_TOKEN")
	if apiToken == "" {
		fmt.Fprintln(os.Stderr, "Error: SLACK_TOKEN environment variable not set")
		fmt.Fprintln(os.Stderr, RootCmd.UsageString())
		os.Exit(1)
	}

	slackChannel, _ := cmd.Flags().GetString("channel")

	client := slackio.NewClient(apiToken)
	defer client.Close()

	reader := slackio.NewReader(client, slackChannel)
	defer reader.Close()

	if _, err := io.Copy(os.Stdout, reader); err != nil {
		panic(err)
	}
}
