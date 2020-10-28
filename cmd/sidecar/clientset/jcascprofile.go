package clientset

import (
	"github.com/jenkinsci/kubernetes-operator/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
)

type JCasCProfileInterface interface {
	List(opts metav1.ListOptions) (*v1alpha2.JCasCProfileList, error)
	Get(name string, options metav1.ListOptions) (*v1alpha2.JCasCProfile, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type jcascprofileClient struct {
	restClient rest.Interface
	namespace  string
}

func (j *jcascprofileClient) List(opts metav1.ListOptions) (*v1alpha2.JCasCProfileList, error) {
	result := v1alpha2.JCasCProfileList{}
	err := j.restClient.
		Get().
		Namespace(j.namespace).
		Resource("jcascprofiles").
		Do().
		Into(&result)

	return &result, err
}

func (j *jcascprofileClient) Get(name string, opts metav1.ListOptions) (*v1alpha2.JCasCProfile, error) {
	result := v1alpha2.JCasCProfile{}
	err := j.restClient.
		Post().
		Namespace(j.namespace).
		Resource("jcascprofiles").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(&result)

	return &result, err
}
func (j *jcascprofileClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return j.restClient.
		Get().
		Namespace(j.namespace).
		Resource("jcascprofiles").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}
