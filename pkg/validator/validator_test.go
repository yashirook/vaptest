package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yashirook/vaptest/pkg/target"
	v1 "k8s.io/api/admissionregistration/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidatePolicy(t *testing.T) {
	testCases := []struct {
		name            string
		policy          *v1.ValidatingAdmissionPolicy
		targetInfoList  target.TargetInfoList
		expectedResults []ValidationResult
		expectedError   string
	}{
		{
			name: "Valid case - Valid object",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "test-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "object.metadata.name.startsWith('test')",
							Message:    "Name must start with 'test'",
						},
					},
				},
			},
			targetInfoList: target.TargetInfoList{
				{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{"name": "test-object"},
					},
					APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-object",
				},
			},
			expectedResults: []ValidationResult{
				{
					PolicyObjectMeta: ObjectMeta{Name: "test-policy"},
					IsValid:          true,
					Message:          "Name must start with 'test'",
					Expression:       "object.metadata.name.startsWith('test')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-object",
					},
				},
			},
		},
		{
			name: "Valid case - Invalid object",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "test-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "object.metadata.name.startsWith('test')",
							Message:    "Name must start with 'test'",
						},
					},
				},
			},
			targetInfoList: target.TargetInfoList{
				{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{"name": "invalid-object"},
					},
					APIGroup: "test.group", APIVersion: "v1", ResourceName: "invalid-object",
				},
			},
			expectedResults: []ValidationResult{
				{
					PolicyObjectMeta: ObjectMeta{Name: "test-policy"},
					IsValid:          false,
					Message:          "Name must start with 'test'",
					Expression:       "object.metadata.name.startsWith('test')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "invalid-object",
					},
				},
			},
		},
		// todo: invalid CEL expression
		// {
		// 	name: "Error case - Invalid CEL expression",
		// 	policy: &v1.ValidatingAdmissionPolicy{
		// 		ObjectMeta: metav1.ObjectMeta{Name: "invalid-policy"},
		// 		Spec: v1.ValidatingAdmissionPolicySpec{
		// 			Validations: []v1.Validation{
		// 				{
		// 					Expression: "invalid.expression",
		// 					Message:    "This expression is invalid",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	targetInfoList: target.TargetInfoList{
		// 		{
		// 			Object: map[string]interface{}{
		// 				"metadata": map[string]interface{}{"name": "test-object"},
		// 			},
		// 			APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-object",
		// 		},
		// 	},
		// 	expectedError: "failed to make AST",
		// },
		{
			name: "Error case - CEL evaluation error",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "error-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "object.nonexistent.field == true",
							Message:    "Accessing a non-existent field",
						},
					},
				},
			},
			targetInfoList: target.TargetInfoList{
				{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{"name": "test-object"},
					},
					APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-object",
				},
			},
			expectedResults: []ValidationResult{},
		},
		{
			name: "ExcludeResourceRules指定時のテスト",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "exclude-resource-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &v1.MatchResources{
						ExcludeResourceRules: []v1.NamedRuleWithOperations{
							{
								RuleWithOperations: v1.RuleWithOperations{
									Rule: v1.Rule{
										APIGroups:   []string{"excluded.group"},
										APIVersions: []string{"v1"},
										Resources:   []string{"excluded-resources"},
									},
								},
								ResourceNames: []string{"excluded-object"},
							},
						},
					},
					Validations: []v1.Validation{
						{
							Expression: "true",
							Message:    "常に有効",
						},
					},
				},
			},
			targetInfoList: target.TargetInfoList{
				{
					Object:       map[string]interface{}{"metadata": map[string]interface{}{"name": "included-object"}},
					APIGroup:     "included.group",
					APIVersion:   "v1",
					Resource:     "included-resources",
					ResourceName: "included-object",
				},
				{
					Object:       map[string]interface{}{"metadata": map[string]interface{}{"name": "excluded-object"}},
					APIGroup:     "excluded.group",
					APIVersion:   "v1",
					Resource:     "excluded-resources",
					ResourceName: "excluded-object",
				},
			},
			expectedResults: []ValidationResult{
				{
					PolicyObjectMeta: ObjectMeta{Name: "exclude-resource-policy"},
					IsValid:          true,
					Message:          "常に有効",
					Expression:       "true",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "included.group",
						ApiGroup:   "v1",
						Name:       "included-object",
					},
				},
			},
		},
		{
			name: "Valid case - Multiple validations",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "multi-validation-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "object.metadata.name.startsWith('test')",
							Message:    "Name must start with 'test'",
						},
						{
							Expression: "object.metadata.name.endsWith('object')",
							Message:    "Name must end with 'object'",
						},
					},
				},
			},
			targetInfoList: target.TargetInfoList{
				{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{"name": "test-valid-object"},
					},
					APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-valid-object",
				},
				{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{"name": "test-object-invalid"},
					},
					APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-object-invalid",
				},
			},
			expectedResults: []ValidationResult{
				{
					PolicyObjectMeta: ObjectMeta{Name: "multi-validation-policy"},
					IsValid:          true,
					Message:          "Name must start with 'test'",
					Expression:       "object.metadata.name.startsWith('test')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-valid-object",
					},
				},
				{
					PolicyObjectMeta: ObjectMeta{Name: "multi-validation-policy"},
					IsValid:          true,
					Message:          "Name must end with 'object'",
					Expression:       "object.metadata.name.endsWith('object')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-valid-object",
					},
				},
				{
					PolicyObjectMeta: ObjectMeta{Name: "multi-validation-policy"},
					IsValid:          true,
					Message:          "Name must start with 'test'",
					Expression:       "object.metadata.name.startsWith('test')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-object-invalid",
					},
				},
				{
					PolicyObjectMeta: ObjectMeta{Name: "multi-validation-policy"},
					IsValid:          false,
					Message:          "Name must end with 'object'",
					Expression:       "object.metadata.name.endsWith('object')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-object-invalid",
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := &Validator{TargetInfoList: tc.targetInfoList}
			results, err := v.validatePolicy(tc.policy)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tc.expectedResults, results)
			}
		})
	}
}
