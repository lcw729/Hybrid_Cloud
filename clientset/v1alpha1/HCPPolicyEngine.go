package v1alpha1

import (
	policyv1alpha1 "Hybrid_Cluster/apis/policy/v1alpha1"
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type HCPPolicyInterface interface {
	List(opts metav1.ListOptions) (*policyv1alpha1.HCPPolicyList, error)
	Get(name string, options metav1.GetOptions) (*policyv1alpha1.HCPPolicy, error)
	Create(deployment *policyv1alpha1.HCPPolicy) (*policyv1alpha1.HCPPolicy, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	// ...
}

type HCPPolicyClient struct {
	restClient rest.Interface
	ns         string
}

func (c *HCPPolicyClient) List(opts metav1.ListOptions) (*policyv1alpha1.HCPPolicyList, error) {
	result := policyv1alpha1.HCPPolicyList{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("hcppolicy").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *HCPPolicyClient) Get(name string, opts metav1.GetOptions) (*policyv1alpha1.HCPPolicy, error) {
	result := policyv1alpha1.HCPPolicy{}
	err := c.restClient.
		Get().
		Namespace(c.ns).
		Resource("hcppolicy").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *HCPPolicyClient) Create(deployment *policyv1alpha1.HCPPolicy) (*policyv1alpha1.HCPPolicy, error) {
	result := policyv1alpha1.HCPPolicy{}
	err := c.restClient.
		Post().
		Namespace(c.ns).
		Resource("hcppolicy").
		Body(deployment).
		Do(context.TODO()).
		Into(&result)

	return &result, err
}

func (c *HCPPolicyClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.
		Get().
		Namespace(c.ns).
		Resource("hcppolicy").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(context.TODO())
}
