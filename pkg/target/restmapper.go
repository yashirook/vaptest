package target

import (
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func createStaticRESTMapper(scheme *runtime.Scheme) meta.RESTMapper {
	// スキームから優先されるグループバージョンを取得
	groupVersions := scheme.PrioritizedVersionsAllGroups()

	// デフォルトのRESTMapperを作成
	mapper := meta.NewDefaultRESTMapper(groupVersions)

	// マッピングを追加
	// Podのマッピング
	podGVK := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}
	podGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "pods"}
	mapper.AddSpecific(podGVK, podGVR, podGVR, meta.RESTScopeNamespace)

	// Serviceのマッピング
	serviceGVK := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"}
	serviceGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "services"}
	mapper.AddSpecific(serviceGVK, serviceGVR, serviceGVR, meta.RESTScopeNamespace)

	// Deploymentのマッピング
	deployGVK := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	deployGVR := schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}
	mapper.AddSpecific(deployGVK, deployGVR, deployGVR, meta.RESTScopeNamespace)

	// 他のリソースも必要に応じて追加
	// ...

	return mapper
}
