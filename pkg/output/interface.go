package output

import (
	"github.com/yashirook/vaptest/pkg/validator"
)

// OutputFormatter は検証結果を出力するためのインターフェースです。
type OutputFormatter interface {
	Output(results []validator.ValidationResult) error
}
