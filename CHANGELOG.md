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

## [v0.1.4] - 2018-06-27
### Changed
- Upgraded internal dependencies to the latest versions.

## [v0.1.3] - 2018-02-24
### Changed
- Package slackio has moved to https://github.com/ahamlinman/slackio
- All dependencies are now version-controlled to support better reproducibility
  of builds

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

## v0.1.0 - 2017-11-24
### Added
- Initial versioned release
- Changelog for all notable modifications going forward

[Unreleased]: https://github.com/ahamlinman/slackbridge/compare/v0.1.4...HEAD
[v0.1.4]: https://github.com/ahamlinman/slackbridge/compare/v0.1.3...v0.1.4
[v0.1.3]: https://github.com/ahamlinman/slackbridge/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/ahamlinman/slackbridge/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/ahamlinman/slackbridge/compare/v0.1.0...v0.1.1
