package cmd

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"go.alexhamlin.co/slackio"

	"go.alexhamlin.co/slackbridge/internal/childproc"
)

var muxCmd = &cobra.Command{
	Use:     "mux [flags] -- program [args]",
	Example: "mux -- ./start-bot.sh -c {{.ChannelID}}",
	Short:   "Exec a separate instance of a program per Slack channel",
	Long: `Mux spawns multiple instances of a provided executable program, where the
standard input, output, and error streams of each instance are connected to a
single Slack channel. The communication semantics are equivalent to those of
Exec mode - this is essentially an Exec mode with automatic support for
multiple channels.

Child processes are spawned on demand when a message is received from a channel
that does not yet have an associated process. The special pattern
"{{.ChannelID}}" will be replaced with the ID of the channel in the arguments
of the spawned process (note that this is a subset of Go's template syntax -
full templating functionality is not supported).

If a child process exits, slackbridge will not automatically respawn it. This
can be used as a basic "filter" to restrict slackbridge to a subset of the
associated user's channels.`,

	Args: cobra.MinimumNArgs(1),
	Run:  runMuxCmd,
}

var channelIDTemplate *regexp.Regexp

func init() {
	RootCmd.AddCommand(muxCmd)

	channelIDTemplate = regexp.MustCompile(`{{\.ChannelID}}`)
}

func runMuxCmd(cmd *cobra.Command, args []string) {
	apiToken := os.Getenv("SLACK_TOKEN")
	if apiToken == "" {
		fmt.Fprintln(os.Stderr, "Error: SLACK_TOKEN environment variable not set")
		fmt.Fprintln(os.Stderr, RootCmd.UsageString())
		os.Exit(1)
	}

	client := slackio.NewClient(apiToken)

	msgs := make(chan slackio.Message)
	client.Subscribe(msgs)

	spawned := make(map[string]bool)

	for msg := range msgs {
		if spawned[msg.ChannelID] {
			continue
		}

		childArgs := make([]string, len(args))
		for i, v := range args {
			childArgs[i] = channelIDTemplate.ReplaceAllString(v, msg.ChannelID)
		}

		reader := slackio.NewReader(&subscriberAt{client, msg.ID}, msg.ChannelID)
		writer := slackio.NewWriter(client, msg.ChannelID, nil)

		// TODO Something other than fire-and-forget...
		childproc.Spawn(childArgs, reader, writer)
		spawned[msg.ChannelID] = true
	}
}

// subscriberAt implements the slackio.ReadClient interface, but starts the
// subscription at a specified message ID using SubscribeAt.
type subscriberAt struct {
	*slackio.Client
	id int
}

func (s *subscriberAt) Subscribe(ch chan<- slackio.Message) error {
	return s.SubscribeAt(s.id, ch)
}
