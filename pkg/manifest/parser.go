package manifest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	scheme = runtime.NewScheme()
	codecs = serializer.NewCodecFactory(scheme)
)

func init() {
	// Register necessary types
	_ = appsv1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
}

// LoadManifests loads Kubernetes manifests from a file or directory
func LoadManifests(path string) ([]runtime.Object, error) {
	var objects []runtime.Object

	info, err := os.ReadDir(path)
	if err != nil {
		return loadManifestFile(path)
	}

	// Process all files in the directory
	for _, entry := range info {
		if entry.IsDir() {
			continue
		}
		filePath := filepath.Join(path, entry.Name())
		objs, err := loadManifestFile(filePath)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}

	fmt.Printf("Object: %v", objects)
	return objects, nil
}

func loadManifestFile(filePath string) ([]runtime.Object, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Split YAML documents
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(data), 1024)
	var objects []runtime.Object

	for {
		var rawObj runtime.RawExtension
		if err := decoder.Decode(&rawObj); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if len(rawObj.Raw) == 0 {
			continue
		}

		obj, _, err := codecs.UniversalDeserializer().Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return nil, err
		}

		objects = append(objects, obj)
	}

	return objects, nil
}
