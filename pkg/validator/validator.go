package validator

import (
	"errors"
	"fmt"

	"github.com/google/cel-go/cel"
	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/cel/environment"
)

type Validator struct {
	TargetObjects  []runtime.Object
	Policies       []*v1.ValidatingAdmissionPolicy
	PolicyBindings []*v1.ValidatingAdmissionPolicyBinding
}

func NewValidator(targetObjects []runtime.Object, policies []*v1.ValidatingAdmissionPolicy, PolicyBindings []*v1.ValidatingAdmissionPolicyBinding) (Validator, error) {
	if len(targetObjects) == 0 {
		return Validator{}, errors.New("target objects is empty")
	}

	if len(policies) == 0 {
		return Validator{}, errors.New("policies is empty")
	}

	if len(PolicyBindings) == 0 {
		return Validator{}, errors.New("policy bindings is empty")
	}

	return Validator{
		TargetObjects:  targetObjects,
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
		return nil, errors.New(fmt.Sprintf("CEL expression parse error: %w\n", issues.Err()))
	}

	prog, err := env.Program(ast)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(fmt.Sprintf("Build CEL Program error: %w\n", err))
	}

	return prog, nil
}

func (v *Validator) validatePolicy(policy *v1.ValidatingAdmissionPolicy) ([]ValidationResult, error) {
	results := make([]ValidationResult, 0)
	for _, validation := range policy.Spec.Validations {

		prog, err := makeCELProgram(&validation)
		if err != nil {
			return results, errors.New(fmt.Sprintf("Failed to make AST: %w\n", err))
		}

		for _, target := range v.TargetObjects {
			objMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(target)
			if err != nil {
				fmt.Println(err)
				return results, err
			}

			activation := map[string]interface{}{
				"object": objMap,
			}

			out, _, err := prog.Eval(activation)
			if err != nil {
				fmt.Printf("eval error: %s\n", err)
				continue
			}
			isValid, ok := out.Value().(bool)
			if !ok {
				return results, errors.New("failed to convert CEL result to bool")
			}

			metadata := objMap["metadata"].(map[string]interface{})
			results = append(results, ValidationResult{
				PolicyObjectMeta: ObjectMeta{
					ApiVersion: policy.APIVersion,
					ApiGroup:   policy.Kind,
					Name:       policy.Name,
				},
				IsValid:    isValid,
				Message:    validation.Message,
				Expression: validation.Expression,
				TargetObjectMeta: ObjectMeta{
					ApiVersion: objMap["apiVersion"].(string),
					ApiGroup:   objMap["kind"].(string),
					Name:       metadata["name"].(string),
					// [TODO] namespaceを返すようにする。namespaceが設定されていない場合のエラーハンドリングが必要
					// Namespace:  metadata["namespace"].(string),
				},
			})
		}
	}
	return results, nil
}
