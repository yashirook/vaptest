package cmd

import (
	"encoding/json"
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

	targetJson, err := json.Marshal(targets)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("target manifests: %s\n", string(targetJson))

	policies, bindings, err := ldr.LoadPolicyFromPaths(policyPaths)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load policy objects: %w", err))
		return
	}

	policyJson, err := json.Marshal(policies)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("policy manifests: %s\n", string(policyJson))
	bindingsJson, err := json.Marshal(bindings)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("policy manifests: %s\n", string(bindingsJson))

	validator := validator.NewValidator(targets, policies, bindings)
	results, err := validator.Validate()
	if err != nil {
		fmt.Printf("validation error: %v", err)
	}

	fmt.Printf("results (len=%d): %v", len(results), results)
}
