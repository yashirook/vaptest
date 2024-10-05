package validator

// func filterTargetObjects(targetObjects []runtime.Object, policy *v1.ValidatingAdmissionPolicy) []runtime.Object {
// 	filteredObjects := make([]runtime.Object, 0)

// 	for _, target := range targetObjects {
// 		info := resourceInfo{
// 			apiGroup:    target.GetObjectKind().GroupVersionKind().Group,
// 			apiVersion:  target.GetObjectKind().GroupVersionKind().Version,
// 			resource:    target.GetObjectKind().GroupVersionKind().Kind,
// 			subResource: target.GetSubresource(),
// 		}
// 	}
// }
