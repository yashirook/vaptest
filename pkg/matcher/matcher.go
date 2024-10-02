package matcher

import (
	"context"
	"errors"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	celtypes "github.com/google/cel-go/common/types"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/admission"
	celplugin "k8s.io/apiserver/pkg/admission/plugin/cel"
	celconfig "k8s.io/apiserver/pkg/apis/cel"

	"k8s.io/apiserver/pkg/admission/plugin/webhook/matchconditions"
)

// matcher evaluates compiled cel expressions and determines if they match the given request or not
type Matcher struct {
	Filter      celplugin.Filter
	MatcherType string
	MatcherKind string
	ObjectName  string
}

// NewMatcher creates a new Matcher with the specified filter, failure policy, matcher type, matcher kind, and object name.
func NewMatcher(filter celplugin.Filter, matcherType, matcherKind, objectName string) Matcher {
	return Matcher{
		Filter:      filter,
		MatcherType: matcherType,
		MatcherKind: matcherKind,
		ObjectName:  objectName,
	}
}

// Match evaluates the given request against the filter and returns true if the request matches the filter, false otherwise.
func (m *Matcher) Match(ctx context.Context, versionedAttr *admission.VersionedAttributes, versionedParams runtime.Object) (matchconditions.MatchResult, error) {
	evalResults, _, err := m.Filter.ForInput(ctx, versionedAttr, celplugin.CreateAdmissionRequest(versionedAttr.Attributes, metav1.GroupVersionResource(versionedAttr.GetResource()), metav1.GroupVersionKind(versionedAttr.VersionedKind)), celplugin.OptionalVariableBindings{
		VersionedParams: versionedParams,
	}, nil, celconfig.RuntimeCELCostBudgetMatchConditions)

	if err != nil {
		return matchconditions.MatchResult{}, err
	}

	errorList := []error{}
	for _, evalResult := range evalResults {
		matchCondition, ok := evalResult.ExpressionAccessor.(*matchconditions.MatchCondition)
		if !ok {
			// This shouldnt happen, but if it does treat same as eval error
			errorList = append(errorList, errors.New(fmt.Sprintf("internal error converting ExpressionAccessor to MatchCondition")))
			continue
		}
		if evalResult.Error != nil {
			errorList = append(errorList, evalResult.Error)
		}
		if evalResult.EvalResult == celtypes.False {
			// If any condition false, skip calling webhook always
			return matchconditions.MatchResult{
				Matches:             false,
				FailedConditionName: matchCondition.Name,
			}, nil
		}
	}
	// if no results eval to false, return matches true with list of any errors encountered
	return matchconditions.MatchResult{
		Matches: true,
	}, nil
}
