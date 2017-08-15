// Package cmd implements the functionality of command slackbridge.
package cmd

import "github.com/spf13/cobra"

// RootCmd is the root of the slackbridge subcommand tree.
var RootCmd = &cobra.Command{
	Use:   "slackbridge",
	Short: "slackbridge connects your command line to Slack",
}
