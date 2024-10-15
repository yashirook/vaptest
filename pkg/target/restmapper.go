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
	addResourceMappings(mapper)

	return mapper
}

func addResourceMappings(mapper *meta.DefaultRESTMapper) {
	// Coreリソース
	addSpecificResource(mapper, "", "v1", "Pod", "pods", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "", "v1", "Service", "services", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "", "v1", "ConfigMap", "configmaps", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "", "v1", "Secret", "secrets", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "", "v1", "Node", "nodes", meta.RESTScopeRoot)
	addSpecificResource(mapper, "", "v1", "Namespace", "namespaces", meta.RESTScopeRoot)
	addSpecificResource(mapper, "", "v1", "PersistentVolume", "persistentvolumes", meta.RESTScopeRoot)
	addSpecificResource(mapper, "", "v1", "PersistentVolumeClaim", "persistentvolumeclaims", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "", "v1", "ServiceAccount", "serviceaccounts", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "", "v1", "ReplicationController", "replicationcontrollers", meta.RESTScopeNamespace)

	// Appsリソース
	addSpecificResource(mapper, "apps", "v1", "Deployment", "deployments", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "apps", "v1", "StatefulSet", "statefulsets", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "apps", "v1", "DaemonSet", "daemonsets", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "apps", "v1", "ReplicaSet", "replicasets", meta.RESTScopeNamespace)

	// Batchリソース
	addSpecificResource(mapper, "batch", "v1", "Job", "jobs", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "batch", "v1", "CronJob", "cronjobs", meta.RESTScopeNamespace)

	// Networkingリソース
	addSpecificResource(mapper, "networking.k8s.io", "v1", "Ingress", "ingresses", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "networking.k8s.io", "v1", "NetworkPolicy", "networkpolicies", meta.RESTScopeNamespace)

	// RBACリソース
	addSpecificResource(mapper, "rbac.authorization.k8s.io", "v1", "Role", "roles", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "rbac.authorization.k8s.io", "v1", "ClusterRole", "clusterroles", meta.RESTScopeRoot)
	addSpecificResource(mapper, "rbac.authorization.k8s.io", "v1", "RoleBinding", "rolebindings", meta.RESTScopeNamespace)
	addSpecificResource(mapper, "rbac.authorization.k8s.io", "v1", "ClusterRoleBinding", "clusterrolebindings", meta.RESTScopeRoot)

	// その他必要なリソースを追加
	// 例:
	// addSpecificResource(mapper, "autoscaling", "v1", "HorizontalPodAutoscaler", "horizontalpodautoscalers", meta.RESTScopeNamespace)
}

func addSpecificResource(mapper *meta.DefaultRESTMapper, group, version, kind, resource string, scope meta.RESTScope) {
	gvk := schema.GroupVersionKind{Group: group, Version: version, Kind: kind}
	gvr := schema.GroupVersionResource{Group: group, Version: version, Resource: resource}
	scopeValue := meta.RESTScopeRoot
	if scope == meta.RESTScopeNamespace {
		scopeValue = meta.RESTScopeNamespace
	}
	mapper.AddSpecific(gvk, gvr, gvr, scopeValue)
}
