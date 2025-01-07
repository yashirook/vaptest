package validator

import (
	"testing"

	"github.com/yashirook/vaptest/pkg/target"
	v1 "k8s.io/api/admissionregistration/v1"
)

func TestMatchesRule(t *testing.T) {
	tests := []struct {
		name       string
		rules      []v1.NamedRuleWithOperations
		targetInfo *target.TargetInfo
		want       bool
	}{
		{
			name:  "Empty rules",
			rules: []v1.NamedRuleWithOperations{},
			targetInfo: &target.TargetInfo{
				TargetIdentifier: target.TargetIdentifier{
					APIGroup:   "apps",
					APIVersion: "v1",
					Resource:   "deployments",
				},
			},
			want: true,
		},
		{
			name: "Exact match",
			rules: []v1.NamedRuleWithOperations{
				{
					RuleWithOperations: v1.RuleWithOperations{
						Rule: v1.Rule{
							APIGroups:   []string{"apps"},
							APIVersions: []string{"v1"},
							Resources:   []string{"deployments"},
						},
					},
				},
			},
			targetInfo: &target.TargetInfo{
				TargetIdentifier: target.TargetIdentifier{
					APIGroup:   "apps",
					APIVersion: "v1",
					Resource:   "deployments",
				},
			},
			want: true,
		},
		{
			name: "Wildcard match",
			rules: []v1.NamedRuleWithOperations{
				{
					RuleWithOperations: v1.RuleWithOperations{
						Rule: v1.Rule{
							APIGroups:   []string{"*"},
							APIVersions: []string{"*"},
							Resources:   []string{"*"},
						},
					},
				},
			},
			targetInfo: &target.TargetInfo{
				TargetIdentifier: target.TargetIdentifier{
					APIGroup:   "apps",
					APIVersion: "v1",
					Resource:   "deployments",
				},
			},
			want: true,
		},
		{
			name: "No match",
			rules: []v1.NamedRuleWithOperations{
				{
					RuleWithOperations: v1.RuleWithOperations{
						Rule: v1.Rule{
							APIGroups:   []string{"core"},
							APIVersions: []string{"v1"},
							Resources:   []string{"pods"},
						},
					},
				},
			},
			targetInfo: &target.TargetInfo{
				TargetIdentifier: target.TargetIdentifier{
					APIGroup:   "apps",
					APIVersion: "v1",
					Resource:   "deployments",
				},
			},
			want: false,
		},
		// Additional test cases can be described here
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchesRule(tt.rules, tt.targetInfo); got != tt.want {
				t.Errorf("matchesRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
