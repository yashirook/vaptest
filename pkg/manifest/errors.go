package manifest

import "fmt"

type FileReadError struct {
	Path string
	Err  error
}

func (e *FileReadError) Error() string {
	return fmt.Sprintf("failed to read file %s: %v", e.Path, e.Err)
}

func (e *FileReadError) Unwrap() error {
	return e.Err
}

type YAMLParseError struct {
	Path string
	Err  error
}

func (e *YAMLParseError) Error() string {
	return fmt.Sprintf("failed to parse YAML file %s: %v", e.Path, e.Err)
}

func (e *YAMLParseError) Unwrap() error {
	return e.Err
}

type UnknownResourceError struct {
	Path    string
	Kind    string
	Version string
}

func (e *UnknownResourceError) Error() string {
	return fmt.Sprintf("unknown resource type in file %s: kind=%s, version=%s", e.Path, e.Kind, e.Version)
}
