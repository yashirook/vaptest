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
	TargetIdentifier
	Object map[string]interface{}
}

type TargetIdentifier struct {
	APIGroup     string `json:"apiGroup"`
	APIVersion   string `json:"apiVersion"`
	Resource     string `json:"resource"`
	SubResource  string `json:"subResource"`
	ResourceName string `json:"resourceName"`
	Namespace    string `json:"namespace"`
}

type TargetInfoList []TargetInfo

func NewTargetInfoList(objects []runtime.Object, scheme *runtime.Scheme) (TargetInfoList, error) {
	results := make([]TargetInfo, 0)
	for _, obj := range objects {
		info, err := NewTargetInfo(obj, scheme)
		if err != nil {
			return nil, err
		}
		results = append(results, *info)
	}
	return results, nil
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

	objMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		fmt.Println(err)
		return &TargetInfo{}, err
	}

	targetInfo := TargetInfo{
		TargetIdentifier: TargetIdentifier{
			APIGroup:     gvk.Group,
			APIVersion:   gvk.Version,
			Resource:     gvr.Resource,
			SubResource:  "", // サブリソースがある場合は設定
			Namespace:    metaObj.GetNamespace(),
			ResourceName: resourceName,
		},
		Object: objMap,
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
