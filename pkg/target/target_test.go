package target

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func TestNewTargetInfo(t *testing.T) {
	testCases := []struct {
		name     string
		obj      runtime.Object
		expected *TargetInfo
		wantErr  bool
	}{
		{
			name: "有効なUnstructuredオブジェクト（Deployment）",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "apps/v1",
					"kind":       "Deployment",
					"metadata": map[string]interface{}{
						"name": "test-deployment",
					},
				},
			},
			expected: &TargetInfo{
				APIGroup:     "apps",
				APIVersion:   "v1",
				Resource:     "deployments",
				ResourceName: "test-deployment",
			},
			wantErr: false,
		},
		{
			name: "有効な構造化オブジェクト（Pod）",
			obj: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "v1",
					Kind:       "Pod",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-pod",
				},
			},
			expected: &TargetInfo{
				APIGroup:     "",
				APIVersion:   "v1",
				Resource:     "pods",
				ResourceName: "test-pod",
			},
			wantErr: false,
		},
		{
			name: "有効な構造化オブジェクト（Deployment）",
			obj: &appsv1.Deployment{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps/v1",
					Kind:       "Deployment",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-structured-deployment",
				},
			},
			expected: &TargetInfo{
				APIGroup:     "apps",
				APIVersion:   "v1",
				Resource:     "deployments",
				ResourceName: "test-structured-deployment",
			},
			wantErr: false,
		},
		{
			name: "未知のリソース種類",
			obj: &unstructured.Unstructured{
				Object: map[string]interface{}{
					"apiVersion": "unknown/v1",
					"kind":       "UnknownResource",
					"metadata": map[string]interface{}{
						"name": "test-unknown",
					},
				},
			},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := NewTargetInfo(tc.obj, scheme.Scheme)

			if (err != nil) != tc.wantErr {
				t.Errorf("NewTargetInfo() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr {
				if result.APIGroup != tc.expected.APIGroup {
					t.Errorf("APIGroup mismatch: got %v, want %v", result.APIGroup, tc.expected.APIGroup)
				}
				if result.APIVersion != tc.expected.APIVersion {
					t.Errorf("APIVersion mismatch: got %v, want %v", result.APIVersion, tc.expected.APIVersion)
				}
				if result.Resource != tc.expected.Resource {
					t.Errorf("Resource mismatch: got %v, want %v", result.Resource, tc.expected.Resource)
				}
				if result.ResourceName != tc.expected.ResourceName {
					t.Errorf("ResourceName mismatch: got %v, want %v", result.ResourceName, tc.expected.ResourceName)
				}
				if result.Object == nil {
					t.Error("Object is nil, but should not be")
				}
			}
		})
	}
}
