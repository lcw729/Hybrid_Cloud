/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/KETI-Hybrid/hcp-pkg/apis/vpa/v1"
	scheme "github.com/KETI-Hybrid/hcp-pkg/client/vpa/v1/clientset/versioned/scheme"
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// VerticalPodAutoscalersGetter has a method to return a VerticalPodAutoscalerInterface.
// A group's client should implement this interface.
type VerticalPodAutoscalersGetter interface {
	VerticalPodAutoscalers(namespace string) VerticalPodAutoscalerInterface
}

// VerticalPodAutoscalerInterface has methods to work with VerticalPodAutoscaler resources.
type VerticalPodAutoscalerInterface interface {
	Create(ctx context.Context, verticalPodAutoscaler *v1.VerticalPodAutoscaler, opts metav1.CreateOptions) (*v1.VerticalPodAutoscaler, error)
	Update(ctx context.Context, verticalPodAutoscaler *v1.VerticalPodAutoscaler, opts metav1.UpdateOptions) (*v1.VerticalPodAutoscaler, error)
	UpdateStatus(ctx context.Context, verticalPodAutoscaler *v1.VerticalPodAutoscaler, opts metav1.UpdateOptions) (*v1.VerticalPodAutoscaler, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.VerticalPodAutoscaler, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.VerticalPodAutoscalerList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VerticalPodAutoscaler, err error)
	VerticalPodAutoscalerExpansion
}

// verticalPodAutoscalers implements VerticalPodAutoscalerInterface
type verticalPodAutoscalers struct {
	client rest.Interface
	ns     string
}

// newVerticalPodAutoscalers returns a VerticalPodAutoscalers
func newVerticalPodAutoscalers(c *AutoscalingV1Client, namespace string) *verticalPodAutoscalers {
	return &verticalPodAutoscalers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the verticalPodAutoscaler, and returns the corresponding verticalPodAutoscaler object, and an error if there is any.
func (c *verticalPodAutoscalers) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.VerticalPodAutoscaler, err error) {
	result = &v1.VerticalPodAutoscaler{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of VerticalPodAutoscalers that match those selectors.
func (c *verticalPodAutoscalers) List(ctx context.Context, opts metav1.ListOptions) (result *v1.VerticalPodAutoscalerList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.VerticalPodAutoscalerList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested verticalPodAutoscalers.
func (c *verticalPodAutoscalers) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a verticalPodAutoscaler and creates it.  Returns the server's representation of the verticalPodAutoscaler, and an error, if there is any.
func (c *verticalPodAutoscalers) Create(ctx context.Context, verticalPodAutoscaler *v1.VerticalPodAutoscaler, opts metav1.CreateOptions) (result *v1.VerticalPodAutoscaler, err error) {
	result = &v1.VerticalPodAutoscaler{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(verticalPodAutoscaler).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a verticalPodAutoscaler and updates it. Returns the server's representation of the verticalPodAutoscaler, and an error, if there is any.
func (c *verticalPodAutoscalers) Update(ctx context.Context, verticalPodAutoscaler *v1.VerticalPodAutoscaler, opts metav1.UpdateOptions) (result *v1.VerticalPodAutoscaler, err error) {
	result = &v1.VerticalPodAutoscaler{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		Name(verticalPodAutoscaler.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(verticalPodAutoscaler).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *verticalPodAutoscalers) UpdateStatus(ctx context.Context, verticalPodAutoscaler *v1.VerticalPodAutoscaler, opts metav1.UpdateOptions) (result *v1.VerticalPodAutoscaler, err error) {
	result = &v1.VerticalPodAutoscaler{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		Name(verticalPodAutoscaler.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(verticalPodAutoscaler).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the verticalPodAutoscaler and deletes it. Returns an error if one occurs.
func (c *verticalPodAutoscalers) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *verticalPodAutoscalers) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched verticalPodAutoscaler.
func (c *verticalPodAutoscalers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.VerticalPodAutoscaler, err error) {
	result = &v1.VerticalPodAutoscaler{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("verticalpodautoscalers").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
