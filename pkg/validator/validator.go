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

func NewValidator(targetObjects []runtime.Object, policies []*v1.ValidatingAdmissionPolicy, PolicyBindings []*v1.ValidatingAdmissionPolicyBinding) Validator {
	return Validator{
		TargetObjects:  targetObjects,
		Policies:       policies,
		PolicyBindings: PolicyBindings,
	}
}

type ValidationResult struct {
	PolicyObjectMeta ObjectMeta
	IsValid          bool
	Message          string
	Expression       string
	TargetObjectMeta ObjectMeta
}

type ObjectMeta struct {
	ApiVersion string
	ApiGroup   string
	Name       string
	Namespace  string
}

func (v *Validator) Validate() ([]ValidationResult, error) {
	results := make([]ValidationResult, 0)
	policy := v.Policies[0]
	prog, err := makeCELProgram(policy)
	if err != nil {
		return results, errors.New(fmt.Sprintf("Failed to make AST: %w\n", err))
	}

	fmt.Printf("target objects length: %d\n", len(v.TargetObjects))

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
			Message:    policy.Spec.Validations[0].Message,
			Expression: policy.Spec.Validations[0].Expression,
			TargetObjectMeta: ObjectMeta{
				ApiVersion: objMap["apiVersion"].(string),
				ApiGroup:   objMap["kind"].(string),
				Name:       metadata["name"].(string),
				// [TODO] namespaceを返すようにする。namespaceが設定されていない場合のエラーハンドリングが必要
				// Namespace:  metadata["namespace"].(string),
			},
		})
	}
	return results, nil
}

func makeCELProgram(policy *v1.ValidatingAdmissionPolicy) (cel.Program, error) {
	celEnv := environment.MustBaseEnvSet(environment.DefaultCompatibilityVersion(), false)
	env := celEnv.NewExpressionsEnv()
	validationRule := policy.Spec.Validations[0]

	ast, issues := env.Parse(validationRule.Expression)
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
	for _, rule := range policy.Spec.Validations {

		prog, err := makeCELProgram(policy)
		if err != nil {
			return results, errors.New(fmt.Sprintf("Failed to make AST: %w\n", err))
		}

		fmt.Printf("target objects length: %d\n", len(v.TargetObjects))

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
				Message:    policy.Spec.Validations[0].Message,
				Expression: policy.Spec.Validations[0].Expression,
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
}
