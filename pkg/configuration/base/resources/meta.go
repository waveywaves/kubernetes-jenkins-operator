package resources

import (
	"fmt"

	"github.com/jenkinsci/kubernetes-operator/api/v1alpha2"
	"github.com/jenkinsci/kubernetes-operator/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewResourceObjectMeta builds ObjectMeta for all Kubernetes resources created by operator
func NewResourceObjectMeta(jenkins *v1alpha2.Jenkins) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      GetResourceName(jenkins),
		Namespace: jenkins.ObjectMeta.Namespace,
		Labels:    BuildResourceLabels(jenkins),
	}
}

// BuildResourceLabels returns labels for all Kubernetes resources created by operator
func BuildResourceLabels(jenkins *v1alpha2.Jenkins) map[string]string {
	return map[string]string{
		constants.LabelAppKey:       constants.LabelAppValue,
		constants.LabelJenkinsCRKey: jenkins.Name,
	}
}

// GetResourceName returns name of Kubernetes resource base on Jenkins CR
func GetResourceName(jenkins *v1alpha2.Jenkins) string {
	return fmt.Sprintf("%s-%s", constants.LabelAppValue, jenkins.ObjectMeta.Name)
}
