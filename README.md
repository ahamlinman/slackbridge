# slackbridge

[![GoDoc](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge?status.svg)](https://godoc.org/gitlab.alexhamlin.co/go/slackbridge)

**slackbridge** connects a Slack channel to the standard I/O streams of an
executable program on a server. This is accomplished by implementing real-time
Slack communication behind Go's [io.ReadWriter] interface, and connecting the
result directly to an [exec.Cmd].

[io.ReadWriter]: https://golang.org/pkg/io/#ReadWriter
[exec.Cmd]: https://golang.org/pkg/os/exec/#Cmd

## Setup

For the time being, slackbridge requires a `config.json` file in the working
directory with the following fields:

* `APIToken`: A Slack API token, likely for a bot user
* `Channel`: The "internal" ID of a single Slack channel (can be obtained from
  the URL path when using Slack in a web browser)

With this created, run slackbridge with the command line of the desired
executable, e.g. `slackbridge cat` to create an echo server. All input and
output is line-buffered.

## Examples

* `slackbridge cat`: echo server
* `slackbridge ed`: collaborative line editing
* `slackbridge sudo bash`: yep, this actually works

## TODOs

* Configuration should not be in a JSON file. My plan is to improve argument
  parsing so that the channel ID can be specified on the command line, possibly
  by name with automatic lookup. The API token will be an environment variable
  so as not to end up in shell history or ps output.
* Every line emitted by the executable becomes an individual Slack message. I
  want to "debounce" these and try to group them together more nicely, so large
  output is transmitted more efficiently.
* Usage documentation and examples should be moved entirely to GoDoc.
