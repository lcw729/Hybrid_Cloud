package v1alpha1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

type KubeFedClusterInterface interface {
	List(opts metav1.ListOptions) (*fedv1b1.KubeFedClusterList, error)
	Get(name string, options metav1.GetOptions) (*fedv1b1.KubeFedCluster, error)
	Create(kubefedcluster *fedv1b1.KubeFedCluster) (*fedv1b1.KubeFedCluster, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type KubeFedClusterClient struct {
	restClient rest.Interface
	ns         string
}

func (c *KubeFedClusterClient) List(opts metav1.ListOptions) (*fedv1b1.KubeFedClusterList, error) {
	result := fedv1b1.KubeFedClusterList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kubefedclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *KubeFedClusterClient) Get(name string, opts metav1.GetOptions) (*fedv1b1.KubeFedCluster, error) {
	result := fedv1b1.KubeFedCluster{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kubefedclusters").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *KubeFedClusterClient) Create(kubefedcluster *fedv1b1.KubeFedCluster) (*fedv1b1.KubeFedCluster, error) {
	result := fedv1b1.KubeFedCluster{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("kubefedclusters").
		Body(kubefedcluster).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *KubeFedClusterClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("kubefedclusters").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
