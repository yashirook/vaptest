package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/yashirook/vaptest/pkg/loader"
)

func validate(cmd *cobra.Command, args []string) {
	scheme := runtime.NewScheme()

	// Register necessary types
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = admissionregistrationv1.AddToScheme(scheme)

	ldr := loader.NewLoader(scheme)
	targets, err := ldr.LoadFromPaths([]string{targetPath})
	if err != nil {
		fmt.Println(fmt.Errorf("failed to load target manifests: %w", err))
		return
	}

	targetJson, err := json.Marshal(targets)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(targetJson))
}
