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
			expectedError:   false,
			expectedResults: []string{"all validation success!"},
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
			expectedValidationErrors: 2,
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
				assert.Contains(t, stderr.String(), tc.expectedResults, "期待するエラーメッセージが含まれていること")
			} else {
				assert.NoError(t, err, "エラーが発生しないことを期待しています")
				for _, expectedResult := range tc.expectedResults {
					assert.Contains(t, stdout.String(), expectedResult, "期待する出力が含まれていること")
				}
				assert.Equal(t, tc.expectedValidationErrors+1, strings.Count(stdout.String(), "\n"), "期待する行数が含まれていること")
			}
		})
	}
}
