package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "0.1.0"
	gitCommit string
	buildDate string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of vaptest",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vaptest version: %s (commit %s, build at %s)\n", version, gitCommit, buildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
