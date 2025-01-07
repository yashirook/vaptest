package output

import (
	"github.com/yashirook/vaptest/pkg/validator"
)

type OutputFormatter interface {
	Output(results []validator.ValidationResult) error
}
