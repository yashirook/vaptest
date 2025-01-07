package validator

import (
	"github.com/yashirook/vaptest/pkg/target"
	v1 "k8s.io/api/admissionregistration/v1"
)

func filterTarget(policy *v1.ValidatingAdmissionPolicy, targetInfoList target.TargetInfoList) (target.TargetInfoList, error) {
	filteredTargets := make(target.TargetInfoList, 0)
	if policy.Spec.MatchConstraints == nil {
		return targetInfoList, nil
	}

	for _, t := range targetInfoList {
		// ExcludeResourceRulesが空でない場合のみチェックを行う
		if len(policy.Spec.MatchConstraints.ExcludeResourceRules) > 0 && matchesExcludeRule(policy.Spec.MatchConstraints.ExcludeResourceRules, &t) {
			continue
		}

		// ResourceRulesが空の場合、デフォルトで全てのリソースにマッチする
		if len(policy.Spec.MatchConstraints.ResourceRules) > 0 && !matchesRule(policy.Spec.MatchConstraints.ResourceRules, &t) {
			continue
		}
		filteredTargets = append(filteredTargets, t)
	}

	return filteredTargets, nil
}
