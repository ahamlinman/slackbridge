/*

Command slackbridge connects Slack channels to system I/O streams using Slack's
real-time messaging API.

Three modes of execution are supported:

The first runs a child process and connects its standard streams to a Slack
channel. Within the child process, the text of individual messages in the
channel is received on stdin. Text emitted on stdout and stderr is sent back to
the channel as individual messages.

The second is similar to the first, but automatically starts a new child
process for each Slack channel from which a message is received.

The third connects to Slack and streams message text to stdout. Input is
ignored.

Communication Model

During its operation, slackbridge needs to convert Slack messages to and from
plain text.

When reading, individual messages are delimited by newlines. Multi-line
messages are equivalent to multiple single-line messages in succession. This
is not configurable.

When writing, lines of output written within a 0.1 second interval are batched
into a single Slack message. This is not configurable through the slackbridge
CLI, though the underlying slackio implementation allows customization of this
"batching" scheme.

Users, reactions, threads, and other Slack features are not represented in any
way. Only the text in the main body of the channel is available. Received
messages are formatted per Slack's "Basic message formatting" as described at
https://api.slack.com/docs/message-formatting. Sent messages should be
formatted in this manner as well. slackbridge does not handle this
automatically.

Usage

Run "slackbridge help" to view full usage information. Before using
slackbridge, the SLACK_TOKEN environment variable must be set to a valid Slack
API token.

Caveats

slackbridge is designed for long-running programs. Extremely short programs
(e.g. a single echo statement in exec mode) are not guaranteed to work as
expected, and issues encountered with slackbridge while running such programs
will not be considered bugs. Excessive runs of short programs with slackbridge
will likely trigger Slack's rate limiting.

*/
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
