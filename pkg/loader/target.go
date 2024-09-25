package loader

import (
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
)

// LoadPolicyFromPaths loads policies and bindings from the specified file paths.
// It returns two slices: one containing ValidatingAdmissionPolicy objects and the other containing ValidatingAdmissionPolicyBinding objects.
// If an error occurs during loading, it returns nil slices and the error.
//
// Parameters:
//   - paths: A slice of strings representing the file paths to load the policies and bindings from.
//
// Returns:
//   - []*admissionregistrationv1.ValidatingAdmissionPolicy: A slice of ValidatingAdmissionPolicy objects.
//   - []*admissionregistrationv1.ValidatingAdmissionPolicyBinding: A slice of ValidatingAdmissionPolicyBinding objects.
//   - error: An error if any occurred during loading, otherwise nil.
func (l *Loader) LoadPolicyFromPaths(paths []string) ([]*admissionregistrationv1.ValidatingAdmissionPolicy, []*admissionregistrationv1.ValidatingAdmissionPolicyBinding, error) {
	objs, err := l.LoadObjectFromPaths(paths)
	if err != nil {
		return nil, nil, err
	}

	var policies []*admissionregistrationv1.ValidatingAdmissionPolicy
	var bindings []*admissionregistrationv1.ValidatingAdmissionPolicyBinding
	for _, obj := range objs {
		switch o := obj.(type) {
		case *admissionregistrationv1.ValidatingAdmissionPolicy:
			policies = append(policies, o)
		case *admissionregistrationv1.ValidatingAdmissionPolicyBinding:
			bindings = append(bindings, o)
		}
	}
	return policies, bindings, nil
}
