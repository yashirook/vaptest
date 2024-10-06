package validator

import "github.com/yashirook/vaptest/pkg/target"

type PolicyIdentifier struct {
	PolicyName string `json:"name"`
}

type ValidationResult struct {
	Target           target.TargetIdentifier `json:"target"`
	Policy           PolicyIdentifier        `json:"policy"`
	Success          bool                    `json:"success"`
	IsValidated      bool                    `json:"isValidated"`
	ValidationErrors []ValidationError       `json:"validationErrors,omitempty"`
}

type ValidationError struct {
	Message string `json:"message"`
	CELExpr string `json:"celExpression"`
}
