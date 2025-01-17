package validator

import (
	"errors"
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/yashirook/vaptest/pkg/target"
	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/cel/environment"
)

type Validator struct {
	TargetInfoList []target.TargetInfo
	Policies       []*v1.ValidatingAdmissionPolicy
	PolicyBindings []*v1.ValidatingAdmissionPolicyBinding
	Scheme         *runtime.Scheme
}

func NewValidator(targets target.TargetInfoList, policies []*v1.ValidatingAdmissionPolicy, PolicyBindings []*v1.ValidatingAdmissionPolicyBinding, scheme *runtime.Scheme) (Validator, error) {
	if len(targets) == 0 {
		return Validator{}, errors.New("target objects is empty")
	}

	if len(policies) == 0 {
		return Validator{}, errors.New("policies is empty")
	}

	for _, policy := range policies {
		if policy.Spec.Validations == nil {
			return Validator{}, fmt.Errorf("policy %s is invalid: validations is empty", policy.Name)
		}
		for _, validation := range policy.Spec.Validations {
			if validation.Expression == "" {
				return Validator{}, fmt.Errorf("policy %s is invalid: validation expression is empty", policy.Name)
			}
		}
	}

	return Validator{
		TargetInfoList: targets,
		Policies:       policies,
		PolicyBindings: PolicyBindings,
	}, nil
}

func (v *Validator) Validate() ([]ValidationResult, error) {
	results := make([]ValidationResult, 0)
	for _, policy := range v.Policies {
		res, err := v.validatePolicy(policy)
		if err != nil {
			return results, err
		}
		results = append(results, res...)
	}
	return results, nil
}

func makeCELProgram(validation *v1.Validation) (cel.Program, error) {
	celEnv := environment.MustBaseEnvSet(environment.DefaultCompatibilityVersion(), false)
	env := celEnv.NewExpressionsEnv()

	ast, issues := env.Parse(validation.Expression)
	if issues != nil && issues.Err() != nil {
		return nil, fmt.Errorf("CEL expression parse error: %w", issues.Err())
	}

	// todo: check implementation
	// _, issues = env.Check(ast)
	// if issues != nil && issues.Err() != nil {
	// 	return nil, fmt.Errorf("CEL expression check error: %w", issues.Err())
	// }

	prog, err := env.Program(ast)
	if err != nil {
		return nil, fmt.Errorf("build CEL Program error: %w", err)
	}

	return prog, nil
}

func (v *Validator) validatePolicy(policy *v1.ValidatingAdmissionPolicy) ([]ValidationResult, error) {
	results := make([]ValidationResult, 0)
	filteredTargets, err := filterTarget(policy, v.TargetInfoList)
	if err != nil {
		return results, fmt.Errorf("failed to filter target: %w", err)
	}
	var isValidated bool = false

	for _, t := range filteredTargets {
		var success bool = true
		validationErrors := make([]ValidationError, 0)
		for _, validation := range policy.Spec.Validations {
			prog, err := makeCELProgram(&validation)
			if err != nil {
				return results, fmt.Errorf("failed to make AST: %w", err)
			}

			activation := map[string]interface{}{
				"object": t.Object,
			}

			out, _, err := prog.Eval(activation)
			if err != nil {
				fmt.Printf("eval error: resource=%s, policy=%s, expression=%s, error=%s\n", t.TargetIdentifier.ResourceName, policy.Name, validation.Expression, err)
				continue
			}
			res, ok := out.Value().(bool)
			if !ok {
				continue
			}

			if !res {
				success = false
				validationErrors = append(validationErrors, ValidationError{
					Message: validation.Message,
					CELExpr: validation.Expression,
				})
			}

			isValidated = true
		}

		if isValidated {
			results = appendResult(results, success, isValidated, policy, t, validationErrors)
		}
	}
	return results, nil
}

func appendResult(results []ValidationResult, success bool, isValidated bool, policy *v1.ValidatingAdmissionPolicy, target target.TargetInfo, validationErrors []ValidationError) []ValidationResult {
	return append(results, ValidationResult{
		Policy: PolicyIdentifier{
			PolicyName: policy.Name,
		},
		Success:          success,
		IsValidated:      isValidated,
		ValidationErrors: validationErrors,
		Target:           target.TargetIdentifier,
	})
}
