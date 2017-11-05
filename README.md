# slackbridge

[![GoDoc](https://godoc.org/github.com/ahamlinman/slackbridge?status.svg)](https://godoc.org/github.com/ahamlinman/slackbridge)

**slackbridge connects your command line to Slack** by transforming messages to
and from lines of text on standard I/O streams. It is powered by the
**slackio** package, which implements real-time Slack communication behind Go's
[io.Reader] and [io.Writer] interfaces.

[io.Reader]: https://golang.org/pkg/io/#Reader
[io.Writer]: https://golang.org/pkg/io/#Writer

## Setup

To use slackbridge, you must obtain a Slack API token and make it available
through the `SLACK_TOKEN` environment variable. Many subcommands also require a
9-character channel ID, which can be obtained from the URL path when viewing
Slack in a browser. This is _not_ the same as the user-visible channel name.

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

1. `git clone https://github.com/ahamlinman/slackbridge.git`
1. `go get github.com/golang/dep/cmd/dep` (If you don't already have [`dep`])
1. `dep ensure`
1. `go run`, `go build`, `go install`, etc.

[`dep`]: https://github.com/golang/dep

## Status and Stability

As of November 2017, the key desired functionalities of command slackbridge and
package slackio have been implemented. Feature development is on an indefinite
hiatus, but maintenance updates (e.g. dependency upgrades) may be made from
time to time.

Although significant API and/or CLI changes are unlikely, a stable version has
not been declared and long-term backwards compatibility is not guaranteed.

## License

MIT (see LICENSE.txt)
