# slackbridge

[![GoDoc](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge?status.svg)](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge)

**slackbridge connects your command line to Slack** by transforming messages to
and from lines of text on standard I/O streams. It is powered by the
**slackio** package, which implements real-time Slack communication behind Go's
[io.Reader] and [io.Writer] interfaces.

[io.Reader]: https://golang.org/pkg/io/#Reader
[io.Writer]: https://golang.org/pkg/io/#Writer

## Setup

To use slackbridge, you must obtain a Slack API token and make it available
through the `SLACK_TOKEN` environment variable. All subcommands also require
a 9-character channel ID, which can be obtained from the URL path when viewing
Slack in a browser. This is _not_ the same as the channel name.

## Usage

slackbridge supports the following capabilities:

* `slackbridge exec`: Run a child process and connect its standard streams to a
  single Slack channel
* `slackbridge mux`: Automatically spawn a child process for each Slack channel
  from which a message is received, with standard streams connected as above
* `slackbridge stream`: Stream messages from a channel to standard output

Run `slackbridge help` for full usage information. Also, see the linked GoDoc
for information on the slackbridge communication model (i.e. how Slack messages
are converted to and from plain text).

## Development

1. `git clone https://gitlab.alexhamlin.co/go/slackbridge.git`
1. `go run`, `go build`, `go install`, etc.

Dependencies are managed with [`dep`] and committed to source control.

[`dep`]: https://github.com/golang/dep

## Status and Stability

As of October 2017, command slackbridge and package slackio are in active
development. APIs are not guaranteed to be stable.

## License

MIT (see LICENSE.txt)
