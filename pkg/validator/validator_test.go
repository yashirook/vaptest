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
					TargetIdentifier: target.TargetIdentifier{
						APIGroup:   "test.group",
						APIVersion: "v1",
						Resource:   "test-object",
					},
				},
			},
			expectedResults: []ValidationResult{
				{
					Policy: PolicyIdentifier{
						PolicyName: "test-policy",
					},
					Success:          true,
					IsValidated:      true,
					ValidationErrors: []ValidationError{},
					Target: target.TargetIdentifier{
						APIGroup:   "test.group",
						APIVersion: "v1",
						Resource:   "test-object",
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
					TargetIdentifier: target.TargetIdentifier{
						APIGroup:   "test.group",
						APIVersion: "v1",
						Resource:   "invalid-object",
					},
				},
			},
			expectedResults: []ValidationResult{
				{
					Policy: PolicyIdentifier{
						PolicyName: "test-policy",
					},
					Success:     false,
					IsValidated: true,
					ValidationErrors: []ValidationError{
						{
							Message: "Name must start with 'test'",
							CELExpr: "object.metadata.name.startsWith('test')",
						},
					},
					Target: target.TargetIdentifier{
						APIGroup:   "test.group",
						APIVersion: "v1",
						Resource:   "invalid-object",
					},
				},
			},
		},
		{
			name: "Error case - Invalid CEL expression",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "invalid-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "invalid",
							Message:    "This expression is invalid",
						},
					},
				},
			},
			targetInfoList: target.TargetInfoList{
				{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{"name": "test-object"},
					},
					TargetIdentifier: target.TargetIdentifier{
						APIGroup:     "test.group",
						APIVersion:   "v1",
						Resource:     "test-object",
						ResourceName: "test-object",
					},
				},
			},
			expectedResults: []ValidationResult{},
		},
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
					TargetIdentifier: target.TargetIdentifier{
						APIGroup:     "test.group",
						APIVersion:   "v1",
						Resource:     "test-object",
						ResourceName: "test-object",
					},
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
					Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "included-object"}},
					TargetIdentifier: target.TargetIdentifier{
						APIGroup:     "included.group",
						APIVersion:   "v1",
						Resource:     "included-resources",
						ResourceName: "included-object",
					},
				},
				{
					Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "excluded-object"}},
					TargetIdentifier: target.TargetIdentifier{
						APIGroup:     "excluded.group",
						APIVersion:   "v1",
						Resource:     "excluded-resources",
						ResourceName: "excluded-object",
					},
				},
			},
			expectedResults: []ValidationResult{
				{
					Policy: PolicyIdentifier{
						PolicyName: "exclude-resource-policy",
					},
					Success:          true,
					IsValidated:      true,
					ValidationErrors: []ValidationError{},
					Target: target.TargetIdentifier{
						APIGroup:     "included.group",
						APIVersion:   "v1",
						Resource:     "included-resources",
						ResourceName: "included-object",
					},
				},
			},
		},
		{
			name: "MatchResourceRules指定時のテスト",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "match-resource-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					MatchConstraints: &v1.MatchResources{
						ResourceRules: []v1.NamedRuleWithOperations{
							{
								RuleWithOperations: v1.RuleWithOperations{
									Rule: v1.Rule{
										APIGroups:   []string{"matched.group"},
										APIVersions: []string{"v1"},
										Resources:   []string{"matched-resources"},
									},
								},
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
					Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "matched-object"}},
					TargetIdentifier: target.TargetIdentifier{
						APIGroup:     "matched.group",
						APIVersion:   "v1",
						Resource:     "matched-resources",
						ResourceName: "matched-object",
					},
				},
				{
					Object: map[string]interface{}{"metadata": map[string]interface{}{"name": "unmatched-object"}},
					TargetIdentifier: target.TargetIdentifier{
						APIGroup:     "unmatched.group",
						APIVersion:   "v1",
						Resource:     "unmatched-resources",
						ResourceName: "unmatched-object",
					},
				},
			},
			expectedResults: []ValidationResult{
				{
					Policy: PolicyIdentifier{
						PolicyName: "match-resource-policy",
					},
					Success:          true,
					IsValidated:      true,
					ValidationErrors: []ValidationError{},
					Target: target.TargetIdentifier{
						APIGroup:     "matched.group",
						APIVersion:   "v1",
						Resource:     "matched-resources",
						ResourceName: "matched-object",
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
					TargetIdentifier: target.TargetIdentifier{
						APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-valid-object",
					},
				},
				{
					Object: map[string]interface{}{
						"metadata": map[string]interface{}{"name": "test-object-invalid"},
					},
					TargetIdentifier: target.TargetIdentifier{
						APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-object-invalid",
					},
				},
			},
			expectedResults: []ValidationResult{
				{
					Policy: PolicyIdentifier{
						PolicyName: "multi-validation-policy",
					},
					Success:          true,
					IsValidated:      true,
					ValidationErrors: []ValidationError{},
					Target: target.TargetIdentifier{
						APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-valid-object",
					},
				},
				{
					Policy: PolicyIdentifier{
						PolicyName: "multi-validation-policy",
					},
					Success:     false,
					IsValidated: true,
					ValidationErrors: []ValidationError{
						{
							Message: "Name must end with 'object'",
							CELExpr: "object.metadata.name.endsWith('object')",
						},
					},
					Target: target.TargetIdentifier{
						APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-object-invalid",
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
