package e2e

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ValidateE2ETest struct {
	name                     string
	targetPaths              []string
	policyPaths              []string
	expectedError            bool
	expectedErrorMessages    []string
	expectedResults          []string
	expectedValidationErrors int
}

func TestValidate(t *testing.T) {
	testCases := []ValidateE2ETest{
		{
			name: "simple_policy_valid",
			targetPaths: []string{
				"testdata/01_simple_policy/valid-target.yaml",
			},
			policyPaths: []string{
				"testdata/01_simple_policy/policy.yaml",
			},
			expectedError:            false,
			expectedValidationErrors: 0,
			expectedResults:          []string{"all validation success!"},
		},
		{
			name: "simple_policy_invalid",
			targetPaths: []string{
				"testdata/01_simple_policy/invalid-target.yaml",
			},
			policyPaths: []string{
				"testdata/01_simple_policy/policy.yaml",
			},
			expectedError:            false,
			expectedResults:          []string{"Deploymentにはラベルが必要です"},
			expectedValidationErrors: 3,
		},
		{
			name: "match_constraints_policy_valid",
			targetPaths: []string{
				"testdata/02_match_constraints_resource_rule_policy/valid-target.yaml",
			},
			policyPaths: []string{
				"testdata/02_match_constraints_resource_rule_policy/policy.yaml",
			},
			expectedError:            false,
			expectedResults:          []string{"all validation success!"},
			expectedValidationErrors: 0,
		},
		{
			name: "match_constraints_policy_invalid",
			targetPaths: []string{
				"testdata/02_match_constraints_resource_rule_policy/invalid-target.yaml",
			},
			policyPaths: []string{
				"testdata/02_match_constraints_resource_rule_policy/policy.yaml",
			},
			expectedError:            false,
			expectedResults:          []string{"Deploymentにはラベルが必要です"},
			expectedValidationErrors: 1,
		},
		{
			name: "match_constraints_exclude_resource_rule_policy_valid",
			targetPaths: []string{
				"testdata/03_match_constraints_exclute_resource_rule_policy/valid-target.yaml",
			},
			policyPaths: []string{
				"testdata/03_match_constraints_exclute_resource_rule_policy/policy.yaml",
			},
			expectedError:   false,
			expectedResults: []string{"all validation success!"},
		},
		{
			name: "match_constraints_exclude_resource_rule_policy_invalid",
			targetPaths: []string{
				"testdata/03_match_constraints_exclute_resource_rule_policy/invalid-target.yaml",
			},
			policyPaths: []string{
				"testdata/03_match_constraints_exclute_resource_rule_policy/policy.yaml",
			},
			expectedError:            false,
			expectedResults:          []string{},
			expectedValidationErrors: 0,
		},
		{
			name: "multiple_target_and_policy_valid",
			targetPaths: []string{
				"testdata/04_multiple_target_and_policy/target1.yaml",
				"testdata/04_multiple_target_and_policy/target2.yaml",
			},
			policyPaths: []string{
				"testdata/04_multiple_target_and_policy/policy1.yaml",
				"testdata/04_multiple_target_and_policy/policy2.yaml",
			},
			expectedError:            false,
			expectedResults:          []string{"Deploymentの名前はappで終わる必要があります", "リソースにはラベルが必要です"},
			expectedValidationErrors: 5,
		},
		// invalid case
		{
			name: "invalid_target",
			targetPaths: []string{
				"testdata/not_exist_target.yaml",
			},
			policyPaths:   []string{},
			expectedError: true,
			expectedErrorMessages: []string{
				"failed to load target manifests",
			},
		},
		{
			name:        "invalid_policy",
			targetPaths: []string{},
			policyPaths: []string{
				"testdata/not_exist_policy.yaml",
			},
			expectedError: true,
			expectedErrorMessages: []string{
				"failed to load policy objects",
			},
		},
		{
			name:        "empty_target",
			targetPaths: []string{},
			policyPaths: []string{
				"testdata/01_simple_policy/policy.yaml",
			},
			expectedError: true,
			expectedErrorMessages: []string{
				"failed to create validator",
			},
		},
		{
			name: "empty_policy",
			targetPaths: []string{
				"testdata/01_simple_policy/valid-target.yaml",
			},
			policyPaths:   []string{},
			expectedError: true,
			expectedErrorMessages: []string{
				"failed to create validator",
			},
		},
		{
			name: "invalid_policy_without_validations",
			targetPaths: []string{
				"testdata/invalid/01_invalid_policy/invalid-target.yaml",
			},
			policyPaths: []string{
				"testdata/invalid/01_invalid_policy/policy-without-validations.yaml",
			},
			expectedError: true,
			expectedErrorMessages: []string{
				"failed to create validator",
			},
		},
		{
			name: "invalid_policy_without_expression",
			targetPaths: []string{
				"testdata/invalid/01_invalid_policy/invalid-target.yaml",
			},
			policyPaths: []string{
				"testdata/invalid/01_invalid_policy/policy-without-expression.yaml",
			},
			expectedError: true,
			expectedErrorMessages: []string{
				"failed to create validator",
			},
		},
		// 対応していないターゲットリソース
		{
			name: "unsupported_target_resource",
			targetPaths: []string{
				"testdata/invalid/02_unsupported_target_resource/target.yaml",
			},
			policyPaths: []string{
				"testdata/invalid/02_unsupported_target_resource/policy.yaml",
			},
			expectedError: true,
			expectedErrorMessages: []string{
				"no kind \"Test\" is registered for version \"test\" in scheme",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{"validate"}
			for _, tp := range tc.targetPaths {
				args = append(args, "--targets", tp)
			}
			for _, pp := range tc.policyPaths {
				args = append(args, "--policies", pp)
			}

			cmd := exec.Command("../../bin/vaptest", args...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr

			err := cmd.Run()

			if tc.expectedError {
				assert.Error(t, err, "エラーが発生することを期待しています")
				for _, expectedError := range tc.expectedErrorMessages {
					assert.Contains(t, stderr.String(), expectedError, "期待するエラーメッセージが含まれていること")
				}
				return
			}

			assert.NoError(t, err, "エラーが発生しないことを期待しています")
			for _, expectedResult := range tc.expectedResults {
				assert.Contains(t, stdout.String(), expectedResult, "期待する出力が含まれていること")
			}
			if tc.expectedValidationErrors > 0 {
				assert.Equal(t, tc.expectedValidationErrors+1, strings.Count(stdout.String(), "\n"), "期待する行数が含まれていること")
			}

		})
	}
}
