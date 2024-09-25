package loader_test

import (
	"path/filepath"
	"testing"

	"github.com/yashirook/vaptest/pkg/loader"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestLoader_LoadPolicyFromPaths(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = admissionregistrationv1.AddToScheme(scheme)

	ldr := loader.NewLoader(scheme)

	tests := []struct {
		name             string
		paths            []string
		wantErr          bool
		expectedPolicies int
		expectedBindings int
	}{
		{
			name:             "ValidSinglePolicy",
			paths:            []string{filepath.Join("testdata", "policy", "valid_single_policy.yaml")},
			wantErr:          false,
			expectedPolicies: 1,
			expectedBindings: 0,
		},
		{
			name:             "ValidSingleBinding",
			paths:            []string{filepath.Join("testdata", "policy", "valid_single_binding.yaml")},
			wantErr:          false,
			expectedPolicies: 0,
			expectedBindings: 1,
		},
		{
			name:             "ValidMultiplePoliciesAndBindings",
			paths:            []string{filepath.Join("testdata", "policy")},
			wantErr:          false,
			expectedPolicies: 1,
			expectedBindings: 1,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			policies, bindings, err := ldr.LoadPolicyFromPaths(tt.paths)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadPolicyFromPaths() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(policies) != tt.expectedPolicies {
				t.Errorf("Expected %d policies, got %d", tt.expectedPolicies, len(policies))
			}
			if len(bindings) != tt.expectedBindings {
				t.Errorf("Expected %d bindings, got %d", tt.expectedBindings, len(bindings))
			}
		})
	}
}
