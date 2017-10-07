package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var muxCmd = &cobra.Command{
	Use:     "mux -- program [args]",
	Example: "mux -- start-child.sh --slack-id {{.ChannelID}}",
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

func init() {
	// RootCmd.AddCommand(muxCmd)
}

func runMuxCmd(cmd *cobra.Command, args []string) {
	fmt.Println("It works!")
}
