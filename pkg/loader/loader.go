package loader

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Loader is a struct for loading Kubernetes resources from YAML files
type Loader struct {
	Scheme *runtime.Scheme
	Codecs serializer.CodecFactory
}

// NewLoader creates a new Loader
func NewLoader(scheme *runtime.Scheme) *Loader {
	return &Loader{
		Scheme: scheme,
		Codecs: serializer.NewCodecFactory(scheme),
	}
}

// LoadObjectFromPaths loads resources from a slice of file or directory paths
func (l *Loader) LoadObjectFromPaths(paths []string) ([]runtime.Object, error) {
	var objects []runtime.Object
	for _, path := range paths {
		objs, err := l.loadFromPath(path)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}
	return objects, nil
}

func (l *Loader) loadFromPath(path string) ([]runtime.Object, error) {
	_, err := os.ReadDir(path)

	if err == nil {
		// ディレクトリの場合
		return l.loadFromDirectory(path)
	}

	// ファイルの場合
	return l.loadFromFile(path)
}

func (l *Loader) loadFromDirectory(dirPath string) ([]runtime.Object, error) {
	var objects []runtime.Object
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, &FileReadError{Path: dirPath, Err: err}
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := filepath.Join(dirPath, file.Name())
		objs, err := l.loadFromFile(filePath)
		if err != nil {
			return nil, err
		}
		objects = append(objects, objs...)
	}
	return objects, nil
}

func (l *Loader) loadFromFile(filePath string) ([]runtime.Object, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, &FileReadError{Path: filePath, Err: err}
	}

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

		obj, gvk, err := l.Codecs.UniversalDeserializer().Decode(rawObj.Raw, nil, nil)
		if err != nil {
			return nil, &DecodeError{Path: filePath, Err: err}
		}

		if _, err := l.Scheme.New(*gvk); err != nil {
			return nil, &UnknownResourceError{Kind: gvk.Kind, Version: gvk.Version}
		}

		objects = append(objects, obj)
	}

	return objects, nil
}
