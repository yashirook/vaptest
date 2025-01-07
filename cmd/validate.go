package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yashirook/vaptest/pkg/loader"
	"github.com/yashirook/vaptest/pkg/output"
	"github.com/yashirook/vaptest/pkg/target"
	"github.com/yashirook/vaptest/pkg/validator"
)

func validate(cmd *cobra.Command, args []string) {

	ldr := loader.NewLoader(scheme)
	targetObjects, err := ldr.LoadObjectFromPaths(targetPaths)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to load target manifests: %w", err))
		os.Exit(1)
	}

	targets, err := target.NewTargetInfoList(targetObjects, scheme)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to create target info list: %w", err))
		os.Exit(1)
	}

	policies, bindings, err := ldr.LoadPolicyFromPaths(policyPaths)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to load policy objects: %w", err))
		os.Exit(1)
	}

	validator, err := validator.NewValidator(targets, policies, bindings, scheme)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to create validator: %w", err))
		os.Exit(1)
	}

	results, err := validator.Validate()
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("validation error: %w", err))
		os.Exit(1)
	}

	formatter := output.NewTableFormatter()
	formatter.Output(results)
}
