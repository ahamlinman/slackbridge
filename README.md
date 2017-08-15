# slackbridge

[![GoDoc](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge?status.svg)](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge)

**slackbridge connects your command line to Slack** by transforming messages to
and from lines of text on standard I/O streams. It is powered by the
**slackio** package, which implements real-time Slack communication behind Go's
[io.ReadWriter] interface.

[io.ReadWriter]: https://golang.org/pkg/io/#ReadWriter

## Setup

To use slackbridge, you must obtain a Slack API token and make it available
through the `SLACK_TOKEN` environment variable. All subcommands also require
a 9-character channel ID, which can be obtained from the URL path when viewing
Slack in a browser. This is _not_ the same as the channel name.

## Usage

slackbridge supports the following capabilities:

* `slackbridge exec`: Run a child process and connect its standard streams to a
  Slack channel
* `slackbridge stream`: Stream messages from a channel to standard output

Run `slackbridge help` for full usage information. Also, see the linked GoDoc
for information on the slackbridge communication model (i.e. how Slack messages
are converted to and from plain text).

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
