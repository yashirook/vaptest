package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashirook/vaptest/pkg/loader"
	"github.com/yashirook/vaptest/pkg/validator"
)

func validate(cmd *cobra.Command, args []string) {

	ldr := loader.NewLoader(scheme)
	targets, err := ldr.LoadObjectFromPaths(targetPaths)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load target manifests: %w", err))
		return
	}

	policies, bindings, err := ldr.LoadPolicyFromPaths(policyPaths)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load policy objects: %w", err))
		return
	}

	validator, err := validator.NewValidator(targets, policies, bindings, scheme)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to create validator: %w", err))
		return
	}

	results, err := validator.Validate()
	if err != nil {
		fmt.Printf("validation error: %v", err)
	}

	fmt.Printf("results (len=%d): %v", len(results), results)
}
