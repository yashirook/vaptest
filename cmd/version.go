package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yashirook/vaptest/pkg/manifest"
)

var (
	version   string
	gitCommit string
	buildDate string

	targetPath string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of vaptest",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("vaptest version: %s (commit %s, build at %s)\n", version, gitCommit, buildDate)
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Kubernetes manifests against ValidationAdmissionPolicies",
	Run: func(cmd *cobra.Command, args []string) {
		// print yaml message
		target, err := manifest.LoadManifests(targetPath)
		if err != nil {
			fmt.Println(err)
			return
		}
		targetJson, err := json.Marshal(target)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(targetJson))
	}}

func init() {
	validateCmd.Flags().StringVarP(&targetPath, "target", "t", "", "Path to the target Kubernetes manifests to validate")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(validateCmd)
}
