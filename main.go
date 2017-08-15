/*

Command slackbridge connects a Slack channel to the standard I/O streams of an
executable program.

I/O Model

Within the child process, the text of individual messages in the channel is
received on stdin. Text emitted on stdout and stderr is sent back to the
channel as individual messages.

When reading, individual messages are delimited by newlines. Multi-line
messages are equivalent to multiple single-line messages in succession.

When writing, lines of output written within a 0.1 second interval are batched
into a single Slack message.

Users, reactions, threads, and other Slack features are not represented in any
way. Only the text in the main body of the channel is available.

Usage

	slackbridge [flags] [program] [arguments...]

The SLACK_TOKEN environment variable must be set to a valid Slack API token.

Flags

	-channel string
		Slack channel ID (required)
	-help
		Show usage information

Examples

	slackbridge -channel C12345678 cat		# simple echo server
	slackbridge -channel DABCDEFGH ed -r		# run the "ed" line editor in restricted mode
	slackbridge -channel DANGEROUS sudo bash	# when you want to live on the edge

Caveats

slackbridge is designed for long-running programs. Extremely short programs
(e.g. a single echo statement) are not guaranteed to work as expected, and
issues encountered with slackbridge while running such programs will not be
considered bugs. Excessive runs of short programs with slackbridge will likely
trigger Slack's rate limiting.

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
