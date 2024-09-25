package cmd

import (
	"os"

	"github.com/spf13/cobra"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	targetPath string
	scheme     = runtime.NewScheme()
)

var rootCmd = &cobra.Command{
	Use:   "vaptest",
	Short: "vaptest is a tool for testing Kubernetes ValidationAdmissionPolicies",
	Long: `vaptest is a CLI tool intended for testing Kubernetes ValidationAdmissionPolicies
and ValidationAdmissionPolicyBindings against actual Kubernetes manifests.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of vaptest",
	Run:   printVersion,
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Kubernetes manifests against ValidationAdmissionPolicies",
	Run:   validate,
}

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Cobra settings
	validateCmd.Flags().StringVarP(&targetPath, "target", "t", "", "Path to the target Kubernetes manifests to validate")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(validateCmd)

	// Register Kubernetes API types
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = admissionregistrationv1.AddToScheme(scheme)
}
