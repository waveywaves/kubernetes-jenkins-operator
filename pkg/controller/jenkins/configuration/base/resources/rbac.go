package resources

import (
	"k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	createVerb = "create"
	deleteVerb = "delete"
	getVerb    = "get"
	listVerb   = "list"
	watchVerb  = "watch"
	patchVerb  = "patch"
	updateVerb = "update"
)

// NewRole returns rbac role for jenkins master
func NewRole(meta metav1.ObjectMeta) *v1.Role {
	return &v1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: meta,
		Rules: []v1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods", "pods/exec", "pods/portforward", "pods/log"},
				Verbs:     []string{createVerb, deleteVerb, getVerb, listVerb, patchVerb, updateVerb, watchVerb},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"secrets", "configmaps"},
				Verbs:     []string{createVerb, deleteVerb, getVerb, listVerb, patchVerb, updateVerb, watchVerb},
			},
			{
				APIGroups: []string{"build.openshift.io", "image.openshift.io"},
				Resources: []string{"*"},
				Verbs:     []string{createVerb, deleteVerb, getVerb, listVerb, patchVerb, updateVerb, watchVerb},
			},
		},
	}
}

// NewRoleBinding returns rbac role binding for jenkins master
func NewRoleBinding(name, namespace, serviceAccountName string, roleRef v1.RoleRef) *v1.RoleBinding {
	return &v1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		RoleRef: roleRef,
		Subjects: []v1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: namespace,
			},
		},
	}
}
