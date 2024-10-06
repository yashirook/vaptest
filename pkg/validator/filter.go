package validator

import (
	"github.com/yashirook/vaptest/pkg/target"
	v1 "k8s.io/api/admissionregistration/v1"
)

func filterTarget(policy *v1.ValidatingAdmissionPolicy, targetInfoList target.TargetInfoList) (target.TargetInfoList, error) {
	filteredTargets := make(target.TargetInfoList, 0)

	for _, t := range targetInfoList {
		if policy.Spec.MatchConstraints != nil {
			excluded := matchesRule(policy.Spec.MatchConstraints.ExcludeResourceRules, &t)
			if excluded {
				continue
			}
			matched := matchesRule(policy.Spec.MatchConstraints.ResourceRules, &t)
			if !matched {
				continue
			}
		}

		filteredTargets = append(filteredTargets, t)
	}

	return filteredTargets, nil
}
