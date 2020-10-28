// Package v1alpha2 contains API Schema definitions for the jenkins.io v1alpha2 API group
// +k8s:deepcopy-gen=package,register
// +groupName=jenkins.io
package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// Kind defines Jenkins CRD kind name
	Kind = "Jenkins"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: "jenkins.io", Version: "v1alpha2"}

// GetObjectKind returns Jenkins object kind
func (in *Jenkins) GetObjectKind() schema.ObjectKind { return in }

// SetGroupVersionKind sets GroupVersionKind
func (in *Jenkins) SetGroupVersionKind(kind schema.GroupVersionKind) {}

// GroupVersionKind returns GroupVersionKind
func (in *Jenkins) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   SchemeGroupVersion.Group,
		Version: SchemeGroupVersion.Version,
		Kind:    Kind,
	}
}

// JenkinsTypeMeta returns Jenkins type meta
func JenkinsTypeMeta() metav1.TypeMeta {
	return metav1.TypeMeta{
		Kind:       Kind,
		APIVersion: SchemeGroupVersion.String(),
	}
}

func init() {
	SchemeBuilder.Register(&Jenkins{}, &JenkinsList{})
	SchemeBuilder.Register(&JenkinsImage{}, &JenkinsImageList{})
	SchemeBuilder.Register(&JCasCProfile{}, &JCasCProfileList{})
}
