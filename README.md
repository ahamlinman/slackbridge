# slackbridge

[![GoDoc](https://godoc.org/go.alexhamlin.co/slackbridge?status.svg)](https://godoc.org/go.alexhamlin.co/slackbridge)
[![Build Status](https://travis-ci.org/ahamlinman/slackbridge.svg?branch=master)](https://travis-ci.org/ahamlinman/slackbridge)

**slackbridge connects your command line to Slack** by transforming messages to
and from lines of text on standard I/O streams. It is powered by the [slackio]
package, which implements real-time Slack communication behind Go's [io.Reader]
and [io.Writer] interfaces.

[slackio]: https://go.alexhamlin.co/slackio
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
1. `go run`, `go build`, `go install`, etc.

Dependencies are managed using the (experimental) [Go modules] feature. When
using Go 1.11+ in module mode, all changes to dependencies will automatically
be tracked. Whenever changes occur, be sure to run `go mod vendor` and commit
changes to source control.

[Go modules]: https://github.com/golang/go/wiki/Modules

## Status and Stability

As of November 2017, the key desired functionalities of slackbridge have been
implemented. My main use case for slackbridge no longer exists, so all
development and maintenance is on an indefinite (and possibly permanent)
hiatus.

Until v1.0.0 is tagged (no guarantees about when or if this will happen), this
project adheres to a scheme based on [Semantic Versioning] as follows:

* MINOR updates could potentially contain breaking changes
* PATCH updates will not contain breaking changes

All notable changes will be documented in CHANGELOG.md.

[Semantic Versioning]: http://semver.org/spec/v2.0.0.html

## License Information

The source code of slackbridge is distributed under the terms of the MIT
License (see LICENSE.txt).

The source code of packages in the `vendor/` directory is distributed under the
terms of those packages' respective licenses. In particular, the following
packages are distributed by HashiCorp, Inc. under the terms of the Mozilla
Public License, version 2.0:

* `vendor/github.com/hashicorp/go-multierror/`
* `vendor/github.com/hashicorp/errwrap/`

A copy of the Mozilla Public License, version 2.0, is available in the LICENSE
file within each of these subdirectories.
