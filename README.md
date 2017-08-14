# slackbridge

[![GoDoc](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge?status.svg)](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge)

**slackbridge** connects a Slack channel to the standard I/O streams of an
executable program on a server. This is accomplished by implementing real-time
Slack communication behind Go's [io.ReadWriter] interface, and connecting the
result directly to an [exec.Cmd].

[io.ReadWriter]: https://golang.org/pkg/io/#ReadWriter
[exec.Cmd]: https://golang.org/pkg/os/exec/#Cmd

## Setup

You must have an API token for slackbridge to connect. The most natural way to
run slackbridge is as a bot user, which can be created through your team's
"Custom Integrations" page. The resulting token must be provided through the
`SLACK_TOKEN` environment variable.

slackbridge filters messages to a single channel using its ID. This is
generally a 9-character identifier that is separate from the channel name, and
can be obtained from the URL path when viewing Slack in a web browser.

See the linked GoDoc above for full usage information and examples.

## License

MIT (see LICENSE.txt)
