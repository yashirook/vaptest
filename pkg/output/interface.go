package output

import (
	"github.com/yashirook/vaptest/pkg/validator"
)

// OutputFormatter は検証結果を出力するためのインターフェースです。
type OutputFormatter interface {
	Format(results []validator.ValidationResult) error
}
