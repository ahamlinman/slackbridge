// Package cmd implements the functionality of command slackbridge.
package cmd

import "github.com/spf13/cobra"

// Version is the current version of slackbridge, which may be injected at
// build time. If it is injected, slackbridge will support an additional
// `--version` flag when built.
var Version string

// RootCmd is the root of the slackbridge subcommand tree.
var RootCmd = &cobra.Command{
	Use:     "slackbridge",
	Short:   "slackbridge connects your command line to Slack",
	Version: Version,
}
