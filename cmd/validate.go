package cmd

import (
	"fmt"

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
		fmt.Println(fmt.Errorf("failed to load target manifests: %w", err))
		return
	}

	targets, err := target.NewTargetInfoList(targetObjects, scheme)

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

	formatter := output.NewDefaultFormatter()
	formatter.Format(results)
}
