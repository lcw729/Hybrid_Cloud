package v1alpha1

import (
	"Hybrid_Cluster/apis"
	resourcev1alpha1 "Hybrid_Cluster/apis/clusterRegister/v1alpha1"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ExampleV1Alpha1Interface interface {
	ClusterRegister(namespace string) ClusterRegisterInterface
	HCPPolicy(namespace string) HCPPolicyInterface
}

type ExampleV1Alpha1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*ExampleV1Alpha1Client, error) {
	apis.AddToScheme(scheme.Scheme)

	config := *c
	config.ContentConfig.GroupVersion = &schema.GroupVersion{Group: resourcev1alpha1.GroupName, Version: resourcev1alpha1.GroupVersion}
	config.APIPath = "/apis"
	//config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &ExampleV1Alpha1Client{restClient: client}, nil
}

func (c *ExampleV1Alpha1Client) ClusterRegister(namespace string) ClusterRegisterInterface {
	return &ClusterRegisterClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}

func (c *ExampleV1Alpha1Client) HCPPolicy(namespace string) HCPPolicyInterface {
	return &HCPPolicyClient{
		restClient: c.restClient,
		ns:         namespace,
	}
}
