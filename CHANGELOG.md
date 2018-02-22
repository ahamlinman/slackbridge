# Changelog

All notable changes to this project will be documented in this file. The format
is based on [Keep a Changelog].

Until v1.0.0 is tagged (no guarantees about when or if this will happen), this
project adheres to a scheme based on [Semantic Versioning] as follows:

* MINOR updates could potentially contain breaking changes
* PATCH updates will not contain breaking changes

[Keep a Changelog]: http://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: http://semver.org/spec/v2.0.0.html

## [Unreleased]

## [v0.1.2] - 2018-02-21
### Changed
- CI builds and tests now use Go 1.10

### Fixed
- Updated Slack library to prevent panicking on certain WebSocket errors (see
  https://github.com/nlopes/slack/issues/260)

## [v0.1.1] - 2018-02-16
### Added
- `licenses` command to print information about the availability of source code
  used to build slackbridge
- `--version` flag to print the current version of slackbridge (when provided
  at build time)

### Security
- Updated Slack library to open WebSocket connections using HTTPS (see
  https://github.com/nlopes/slack/pull/208)

## [v0.1.0] - 2017-11-24
### Added
- Initial versioned release
- Changelog for all notable modifications going forward

[Unreleased]: https://github.com/ahamlinman/slackbridge/compare/v0.1.2...HEAD
[v0.1.2]: https://github.com/ahamlinman/slackbridge/tree/v0.1.2
[v0.1.1]: https://github.com/ahamlinman/slackbridge/tree/v0.1.1
[v0.1.0]: https://github.com/ahamlinman/slackbridge/tree/v0.1.0
