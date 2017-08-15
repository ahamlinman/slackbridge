# slackbridge

[![GoDoc](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge?status.svg)](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge)

**slackbridge connects your command line to Slack** by transforming messages to
and from lines of text on standard I/O streams. It is powered by the
**slackio** package, which implements real-time Slack communication behind Go's
[io.ReadWriter] interface.

[io.ReadWriter]: https://golang.org/pkg/io/#ReadWriter

## Setup

You must have an API token for slackbridge to connect. The most natural way to
run slackbridge is as a bot user, which can be created through your team's
"Custom Integrations" page. The resulting token must be provided through the
`SLACK_TOKEN` environment variable.

When connecting to a single channel, slackbridge requires a channel ID. This is
generally a 9-character identifier that is **not** the channel name. It can be
obtained from the URL path when viewing Slack in a web browser.

## Usage

slackbridge's only current execution mode is to run a child process and connect
a Slack channel to its standard input and output streams. Other modes of
execution will be supported in the future.

Run `slackbridge help` for full usage information. Also, see the linked GoDoc
above for more information about how communication via slackbridge works.

## Development

1. `git clone https://gitlab.alexhamlin.co/go/slackbridge.git`
1. `go get -u github.com/golang/dep/cmd/dep`
1. `dep ensure`
1. `go run`, `go build`, `go install`, etc.

## Status and Stability

As of August 2017, command slackbridge and package slackio are in active
development. APIs are not guaranteed to be stable.

## License

MIT (see LICENSE.txt)
