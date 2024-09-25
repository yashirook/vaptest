package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/yashirook/vaptest/pkg/loader"
)

func validate(cmd *cobra.Command, args []string) {

	ldr := loader.NewLoader(scheme)
	targets, err := ldr.LoadFromPaths(targetPaths)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load target manifests: %w", err))
		return
	}

	targetJson, err := json.Marshal(targets)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("target manifests: %s", string(targetJson))

	policies, err := ldr.LoadFromPaths(policyPaths)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load policy manifests: %w", err))
		return
	}

	policyJson, err := json.Marshal(policies)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("policy manifests: %s", string(policyJson))
}
