/*

Command slackbridge connects a Slack channel to the standard I/O streams of an
executable program.

I/O Model

Within the child process, the text of individual messages in the channel is
received on stdin. Text emitted on stdout and stderr is sent back to the
channel as individual messages.

When reading, individual messages are delimited by newlines. Multi-line
messages are equivalent to multiple single-line messages in succession.

When writing, each line is sent as a separate message. No "batching" of
multiple writes is performed.

Users, reactions, threads, and other Slack features are not represented in any
way. Only the text in the main body of the channel is available.

Usage

	slackbridge [executable] [arguments ...]

Examples

	slackbridge cat		# simple echo server
	slackbridge ed -r	# run the "ed" line editor in restricted mode

*/
package main
