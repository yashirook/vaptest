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
			name: "正常なケース - 有効なオブジェクト",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "test-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "object.metadata.name.startsWith('test')",
							Message:    "名前は'test'で始まる必要があります",
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
					Message:          "名前は'test'で始まる必要があります",
					Expression:       "object.metadata.name.startsWith('test')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-object",
					},
				},
			},
		},
		{
			name: "正常なケース - 無効なオブジェクト",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "test-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "object.metadata.name.startsWith('test')",
							Message:    "名前は'test'で始まる必要があります",
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
					Message:          "名前は'test'で始まる必要があります",
					Expression:       "object.metadata.name.startsWith('test')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "invalid-object",
					},
				},
			},
		},
		// todo: invalid CEL expression
		// {
		// 	name: "異常系 - 無効なCEL式",
		// 	policy: &v1.ValidatingAdmissionPolicy{
		// 		ObjectMeta: metav1.ObjectMeta{Name: "invalid-policy"},
		// 		Spec: v1.ValidatingAdmissionPolicySpec{
		// 			Validations: []v1.Validation{
		// 				{
		// 					Expression: "invalid.expression",
		// 					Message:    "この式は無効です",
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
			name: "異常系 - CEL評価エラー",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "error-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "object.nonexistent.field == true",
							Message:    "存在しないフィールドにアクセスしています",
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
			name: "正常なケース - 複数の検証",
			policy: &v1.ValidatingAdmissionPolicy{
				ObjectMeta: metav1.ObjectMeta{Name: "multi-validation-policy"},
				Spec: v1.ValidatingAdmissionPolicySpec{
					Validations: []v1.Validation{
						{
							Expression: "object.metadata.name.startsWith('test')",
							Message:    "名前は'test'で始まる必要があります",
						},
						{
							Expression: "object.metadata.name.endsWith('object')",
							Message:    "名前は'object'で終わる必要があります",
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
						"metadata": map[string]interface{}{"name": "test-valid2-object"},
					},
					APIGroup: "test.group", APIVersion: "v1", ResourceName: "test-valid2-object",
				},
			},
			expectedResults: []ValidationResult{
				{
					PolicyObjectMeta: ObjectMeta{Name: "multi-validation-policy"},
					IsValid:          true,
					Message:          "名前は'test'で始まる必要があります",
					Expression:       "object.metadata.name.startsWith('test')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-valid-object",
					},
				},
				{
					PolicyObjectMeta: ObjectMeta{Name: "multi-validation-policy"},
					IsValid:          true,
					Message:          "名前は'object'で終わる必要があります",
					Expression:       "object.metadata.name.endsWith('object')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-valid-object",
					},
				},
				{
					PolicyObjectMeta: ObjectMeta{Name: "multi-validation-policy"},
					IsValid:          true,
					Message:          "名前は'test'で始まる必要があります",
					Expression:       "object.metadata.name.startsWith('test')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-valid2-object",
					},
				},
				{
					PolicyObjectMeta: ObjectMeta{Name: "multi-validation-policy"},
					IsValid:          true,
					Message:          "名前は'object'で終わる必要があります",
					Expression:       "object.metadata.name.endsWith('object')",
					TargetObjectMeta: ObjectMeta{
						ApiVersion: "test.group", ApiGroup: "v1", Name: "test-valid2-object",
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
