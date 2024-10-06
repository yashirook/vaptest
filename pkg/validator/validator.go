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
	filteredTargets := make(target.TargetInfoList, 0)

	for _, t := range v.TargetInfoList {
		if policy.Spec.MatchConstraints != nil {
			matched := matchesRule(policy.Spec.MatchConstraints.ResourceRules, &t)
			if !matched {
				continue
			}
		}

		filteredTargets = append(filteredTargets, t)
	}

	for _, validation := range policy.Spec.Validations {
		prog, err := makeCELProgram(&validation)
		if err != nil {
			return results, fmt.Errorf("failed to make AST: %w", err)
		}

		for _, t := range filteredTargets {
			activation := map[string]interface{}{
				"object": t.Object,
			}

			out, _, err := prog.Eval(activation)
			if err != nil {
				fmt.Printf("eval error: %s\n", err)
				continue
			}
			isValid, ok := out.Value().(bool)
			if !ok {
				return results, fmt.Errorf("failed to convert CEL result to bool")
			}

			results = appendResult(results, isValid, policy, t, validation)
		}
	}
	return results, nil
}

func appendResult(results []ValidationResult, isValid bool, policy *v1.ValidatingAdmissionPolicy, target target.TargetInfo, validation v1.Validation) []ValidationResult {
	return append(results, ValidationResult{
		PolicyObjectMeta: ObjectMeta{
			ApiVersion: policy.APIVersion,
			ApiGroup:   policy.Kind,
			Name:       policy.Name,
		},
		IsValid:    isValid,
		Message:    validation.Message,
		Expression: validation.Expression,
		TargetObjectMeta: ObjectMeta{
			ApiVersion: target.APIGroup,
			ApiGroup:   target.APIVersion,
			Name:       target.ResourceName,
			// [TODO] namespaceを返すようにする。namespaceが設定されていない場合のエラーハンドリングが必要
			// Namespace:  metadata["namespace"].(string),
		},
	})
}
