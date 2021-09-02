package v1alpha1

import (
	resourcev1alpha1 "Hybrid_Cluster/apis/clusterRegister/v1alpha1"
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ClusterRegisterInterface interface {
	List(opts metav1.ListOptions) (*resourcev1alpha1.ClusterRegisterList, error)
	Get(name string, options metav1.GetOptions) (*resourcev1alpha1.ClusterRegister, error)
	Create(deployment *resourcev1alpha1.ClusterRegister) (*resourcev1alpha1.ClusterRegister, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type ClusterRegisterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *ClusterRegisterClient) List(opts metav1.ListOptions) (*resourcev1alpha1.ClusterRegisterList, error) {
	result := resourcev1alpha1.ClusterRegisterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("ClusterRegister").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ClusterRegisterClient) Get(name string, opts metav1.GetOptions) (*resourcev1alpha1.ClusterRegister, error) {
	result := resourcev1alpha1.ClusterRegister{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("ClusterRegister").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ClusterRegisterClient) Create(deployment *resourcev1alpha1.ClusterRegister) (*resourcev1alpha1.ClusterRegister, error) {
	result := resourcev1alpha1.ClusterRegister{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("ClusterRegister").
		Body(deployment).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *ClusterRegisterClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("ClusterRegister").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
