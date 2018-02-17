package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var licensesCmd = &cobra.Command{
	Use:   "licenses",
	Short: "Learn how you can obtain the source code slackbridge was built with",
	Run:   runLicensesCmd,
}

func init() {
	RootCmd.AddCommand(licensesCmd)
}

func runLicensesCmd(_ *cobra.Command, _ []string) {
	fmt.Println(`
The source code for slackbridge is available to you under the terms of the MIT
License, and can be found at https://github.com/ahamlinman/slackbridge.

This release of slackbridge also incorporates code from the go-multierror and
errwrap libraries distributed by HashiCorp, Inc. under the terms of the Mozilla
Public License, version 2.0. The source code of these libraries is available to
you at the following locations:

  - go-multierror: https://github.com/hashicorp/go-multierror
  - errwrap: https://github.com/hashicorp/errwrap
`)
}
