package validator

import (
	"strings"

	"github.com/yashirook/vaptest/pkg/target"
	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

type resourceInfo struct {
	apiGroup     string
	apiVersion   string
	resource     string
	subResource  string
	resourceName string
	operation    string
}

func getObjectMeta(obj runtime.Object) metav1.Object {
	if accessor, err := meta.Accessor(obj); err == nil {
		return accessor
	}
	return nil
}

func getGVK(obj runtime.Object, scheme *runtime.Scheme) (*schema.GroupVersionKind, error) {
	gvk, err := apiutil.GVKForObject(obj, scheme)
	if err != nil {
		return nil, err
	}
	return &gvk, nil
}

func matchesRule(rules []v1.NamedRuleWithOperations, targetInfo *target.TargetInfo) bool {
	if len(rules) == 0 || rules == nil {
		return true
	}

	for _, rule := range rules {
		if !matchesString(rule.APIGroups, targetInfo.APIGroup) {
			return false
		}
		if !matchesString(rule.APIVersions, targetInfo.APIVersion) {
			return false
		}
		if !matchesResource(rule.Resources, targetInfo.Resource, targetInfo.SubResource) {
			return false
		}
		if len(rule.ResourceNames) > 0 && !matchesString(rule.ResourceNames, targetInfo.ResourceName) {
			return false
		}
		// OperationPolicy is not supported.
	}

	return true
}

func matchesString(patterns []string, value string) bool {
	if len(patterns) == 0 {
		// パターンが指定されていない場合はマッチ
		return true
	}
	for _, pattern := range patterns {
		if pattern == "*" || pattern == value {
			return true
		}
	}
	return false
}

func matchesResource(patterns []string, resource string, subResource string) bool {
	fullResource := resource
	if subResource != "" {
		fullResource = resource + "/" + subResource
	}
	if len(patterns) == 0 {
		// パターンが指定されていない場合はマッチ
		return true
	}
	for _, pattern := range patterns {
		if pattern == "*" || pattern == resource || pattern == fullResource {
			return true
		}
		// ワイルドカード付きのパターンを処理（例: "pods/*"）
		if strings.HasSuffix(pattern, "/*") {
			baseResource := strings.TrimSuffix(pattern, "/*")
			if baseResource == resource {
				return true
			}
		}
	}
	return false
}
