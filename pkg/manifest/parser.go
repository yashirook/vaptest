package manifest

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Process all dir in the directory
	info, err := os.ReadDir(path)
	if err != nil {
		objs, err := loadManifestFile(path)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
		return objects, nil
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

	return objects, nil
}

func loadManifestFile(filePath string) ([]runtime.Object, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &FileReadError{Path: filePath, Err: err}
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
			return nil, &YAMLParseError{Path: filePath, Err: err}
		}

		if len(rawObj.Raw) == 0 {
			continue
		}

		obj, gvk, err := codecs.UniversalDeserializer().Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return nil, &UnknownResourceError{
				Path:    filePath,
				Kind:    gvk.Kind,
				Version: gvk.Version,
			}
		}

		objects = append(objects, obj)
	}

	return objects, nil
}

// ObjectMeta retrieves the ObjectMeta from a runtime.Object
func ObjectMeta(obj runtime.Object) (metav1.Object, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	return accessor, nil
}
