package loader

import "fmt"

// FileReadError represents an error that occurs when reading a file
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

// YAMLParseError represents an error that occurs when parsing YAML
type YAMLParseError struct {
	Path string
	Err  error
}

func (e *YAMLParseError) Error() string {
	return fmt.Sprintf("failed to parse YAML in file %s: %v", e.Path, e.Err)
}

func (e *YAMLParseError) Unwrap() error {
	return e.Err
}

// DecodeError represents an error that occurs when decoding an object
type DecodeError struct {
	Path string
	Err  error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("failed to decode object in file %s: %v", e.Path, e.Err)
}

func (e *DecodeError) Unwrap() error {
	return e.Err
}

// UnknownResourceError は未知のリソースが存在する場合のエラーを表します
type UnknownResourceError struct {
	Kind    string
	Version string
}

func (e *UnknownResourceError) Error() string {
	return fmt.Sprintf("unknown resource type: kind=%s, version=%s", e.Kind, e.Version)
}
