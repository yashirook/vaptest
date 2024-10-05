package target

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
)

type TargetInfo struct {
	APIGroup     string
	APIVersion   string
	Resource     string
	SubResource  string
	ResourceName string
	Object       map[string]interface{}
}

func NewTargetInfo(obj runtime.Object, scheme *runtime.Scheme) (*TargetInfo, error) {
	metaObj, err := getObjectMeta(obj)
	if err != nil {
		return &TargetInfo{}, err
	}

	gvk, err := getGroupVersionKind(obj)
	if err != nil {
		return &TargetInfo{}, err
	}

	// 静的なRESTMapperを作成
	mapper := createStaticRESTMapper(scheme)

	gvr, err := getGroupVersionResource(gvk, mapper)
	if err != nil {
		return &TargetInfo{}, err
	}

	resourceName := metaObj.GetName()

	targetInfo := TargetInfo{
		APIGroup:     gvk.Group,
		APIVersion:   gvk.Version,
		Resource:     gvr.Resource,
		SubResource:  "", // サブリソースがある場合は設定
		ResourceName: resourceName,
	}

	return &targetInfo, nil
}

func getObjectMeta(obj runtime.Object) (metav1.Object, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}
	return accessor, nil
}

func getGroupVersionKind(obj runtime.Object) (schema.GroupVersionKind, error) {
	gvk := obj.GetObjectKind().GroupVersionKind()
	if gvk.Empty() {
		// GVKが設定されていない場合、スキームから取得
		gvks, _, err := scheme.Scheme.ObjectKinds(obj)
		if err != nil {
			return schema.GroupVersionKind{}, err
		}
		if len(gvks) == 0 {
			return schema.GroupVersionKind{}, fmt.Errorf("GVK not found for object")
		}
		gvk = gvks[0]
	}
	return gvk, nil
}

func getGroupVersionResource(gvk schema.GroupVersionKind, mapper meta.RESTMapper) (schema.GroupVersionResource, error) {
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return schema.GroupVersionResource{}, err
	}
	return mapping.Resource, nil
}
