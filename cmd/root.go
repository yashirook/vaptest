package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaptest",
	Short: "vaptest is a tool for testing Kubernetes ValidationAdmissionPolicies",
	Long: `vaptest is a CLI tool intended for testing Kubernetes ValidationAdmissionPolicies
and ValidationAdmissionPolicyBindings against actual Kubernetes manifests.`,
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
