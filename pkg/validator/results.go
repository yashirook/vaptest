package validator

type ValidationResult struct {
	PolicyObjectMeta ObjectMeta
	IsValid          bool
	Message          string
	Expression       string
	TargetObjectMeta ObjectMeta
}

type ObjectMeta struct {
	ApiVersion string
	ApiGroup   string
	Name       string
	Namespace  string
}
