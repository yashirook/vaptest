package loader_test

import (
	"path/filepath"
	"testing"

	"github.com/yashirook/vaptest/pkg/loader"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestLoader_LoadFromPaths(t *testing.T) {
	scheme := runtime.NewScheme()
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	_ = admissionregistrationv1.AddToScheme(scheme)

	ldr := loader.NewLoader(scheme)

	tests := []struct {
		name     string
		paths    []string
		wantErr  bool
		expected int
	}{
		{
			name:     "ValidSingleManifest",
			paths:    []string{filepath.Join("testdata", "valid_single_manifest.yaml")},
			wantErr:  false,
			expected: 1,
		},
		{
			name:     "ValidMultipleManifestsInSingleFile",
			paths:    []string{filepath.Join("testdata", "valid_multiple_manifests.yaml")},
			wantErr:  false,
			expected: 2,
		},
		{
			name:    "InvalidYAML",
			paths:   []string{filepath.Join("testdata", "invalid_yaml.yaml")},
			wantErr: true,
		},
		{
			name:    "UnknownResource",
			paths:   []string{filepath.Join("testdata", "unknown_resource.yaml")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			objs, err := ldr.LoadFromPaths(tt.paths)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadFromPaths() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(objs) != tt.expected {
				t.Errorf("Expected %d objects, got %d", tt.expected, len(objs))
			}
		})
	}
}