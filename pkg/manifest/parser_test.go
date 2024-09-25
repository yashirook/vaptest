package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/yashirook/vaptest/pkg/manifest"
)

func TestLoadManifests(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		wantErr       bool
		expectedCount int
		expectedKinds []string
		expectedNames []string
	}{
		{
			name:          "ValidSingleManifest",
			path:          filepath.Join("testdata", "valid_single_manifest.yaml"),
			wantErr:       false,
			expectedCount: 1,
			expectedKinds: []string{"Deployment"},
			expectedNames: []string{"nginx-deployment"},
		},
		{
			name:          "ValidMultipleManifestsInSingleFile",
			path:          filepath.Join("testdata", "valid_multiple_manifests.yaml"),
			wantErr:       false,
			expectedCount: 2,
			expectedKinds: []string{"Namespace", "Deployment"},
			expectedNames: []string{"test-namespace", "test-deployment"},
		},
		{
			name:          "ValidMultipleFilesInDirectory",
			path:          filepath.Join("testdata", "multiple_files"),
			wantErr:       false,
			expectedCount: 2,
			expectedKinds: []string{"Deployment", "Service"},
			expectedNames: []string{"nginx-deployment", "nginx-service"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifests, err := manifest.LoadManifests(tt.path)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadManifests() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if len(manifests) != tt.expectedCount {
				t.Errorf("Expected %d manifests, got %d", tt.expectedCount, len(manifests))
			}

			for i, obj := range manifests {
				gvk := obj.GetObjectKind().GroupVersionKind()
				if gvk.Kind != tt.expectedKinds[i] {
					t.Errorf("Expected kind %q, got %q", tt.expectedKinds[i], gvk.Kind)
				}

				accessor, err := manifest.ObjectMeta(obj)
				if err != nil {
					t.Errorf("Failed to get object meta: %v", err)
				}

				if accessor.GetName() != tt.expectedNames[i] {
					t.Errorf("Expected name %q, got %q", tt.expectedNames[i], accessor.GetName())
				}
			}
		})
	}
}
