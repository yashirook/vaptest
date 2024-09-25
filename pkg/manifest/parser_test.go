package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/yashirook/vaptest/pkg/manifest"
	appsv1 "k8s.io/api/apps/v1"
)

func TestLoadManifests(t *testing.T) {
	testManifestDir := filepath.Join("testdata", "manifests")

	manifests, err := manifest.LoadManifests(testManifestDir)
	if err != nil {
		t.Fatalf("Failed to load manifests: %v", err)
	}

	expectedCount := 1
	if len(manifests) != expectedCount {
		t.Errorf("Expected %d manifests, got %d", expectedCount, len(manifests))
	}

	obj := manifests[0]
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		t.Fatalf("Expected a Deployment, got %T", obj)
	}

	expectedName := "nginx-deployment"
	if deployment.Name != expectedName {
		t.Errorf("Expected Deployment name %q, got %q", expectedName, deployment.Name)
	}
}
