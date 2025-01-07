package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   string
	gitCommit string
	buildDate string
)

func printVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("vaptest version: %s (commit %s, build at %s)\n", version, gitCommit, buildDate)
}
