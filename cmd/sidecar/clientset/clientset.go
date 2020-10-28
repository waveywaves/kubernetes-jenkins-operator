package clientset

import (
	"github.com/jenkinsci/kubernetes-operator/api/v1alpha2"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

type V1Alpha2Interface interface {
	JCasCProfiles(namespace string) JCasCProfileInterface
}

type V1Alpha2Client struct {
	RESTClient rest.Interface
}

func GetInClusterConfig() *rest.Config {
	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		klog.Fatal(err)
	}

	return config
}

func NewForConfig(config *rest.Config) (*V1Alpha2Client, error) {
	jenkinsResourcesConfig := *config
	jenkinsResourcesConfig.ContentConfig.GroupVersion = &v1alpha2.SchemeGroupVersion
	jenkinsResourcesConfig.APIPath = "/apis"
	jenkinsResourcesConfig.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	jenkinsResourcesConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&jenkinsResourcesConfig)
	if err != nil {
		klog.Fatal(err)
	}

	return &V1Alpha2Client{RESTClient: client}, nil
}

func (c *V1Alpha2Client) JCasCProfiles(namespace string) JCasCProfileInterface {
	return &jcascprofileClient{
		restClient: c.RESTClient,
		namespace:  namespace,
	}
}
