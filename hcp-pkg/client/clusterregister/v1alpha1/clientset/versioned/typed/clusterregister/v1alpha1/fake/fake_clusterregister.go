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

package fake

import (
	v1alpha1 "hcp-pkg/apis/clusterregister/v1alpha1"
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeClusterRegisters implements ClusterRegisterInterface
type FakeClusterRegisters struct {
	Fake *FakeHcpV1alpha1
	ns   string
}

var clusterregistersResource = schema.GroupVersionResource{Group: "hcp.crd.com", Version: "v1alpha1", Resource: "clusterregisters"}

var clusterregistersKind = schema.GroupVersionKind{Group: "hcp.crd.com", Version: "v1alpha1", Kind: "ClusterRegister"}

// Get takes name of the clusterRegister, and returns the corresponding clusterRegister object, and an error if there is any.
func (c *FakeClusterRegisters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ClusterRegister, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(clusterregistersResource, c.ns, name), &v1alpha1.ClusterRegister{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegister), err
}

// List takes label and field selectors, and returns the list of ClusterRegisters that match those selectors.
func (c *FakeClusterRegisters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ClusterRegisterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(clusterregistersResource, clusterregistersKind, c.ns, opts), &v1alpha1.ClusterRegisterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ClusterRegisterList{TypeMeta: obj.(*v1alpha1.ClusterRegisterList).TypeMeta}
	for _, item := range obj.(*v1alpha1.ClusterRegisterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested clusterRegisters.
func (c *FakeClusterRegisters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(clusterregistersResource, c.ns, opts))

}

// Create takes the representation of a clusterRegister and creates it.  Returns the server's representation of the clusterRegister, and an error, if there is any.
func (c *FakeClusterRegisters) Create(ctx context.Context, clusterRegister *v1alpha1.ClusterRegister, opts v1.CreateOptions) (result *v1alpha1.ClusterRegister, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(clusterregistersResource, c.ns, clusterRegister), &v1alpha1.ClusterRegister{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegister), err
}

// Update takes the representation of a clusterRegister and updates it. Returns the server's representation of the clusterRegister, and an error, if there is any.
func (c *FakeClusterRegisters) Update(ctx context.Context, clusterRegister *v1alpha1.ClusterRegister, opts v1.UpdateOptions) (result *v1alpha1.ClusterRegister, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(clusterregistersResource, c.ns, clusterRegister), &v1alpha1.ClusterRegister{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegister), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeClusterRegisters) UpdateStatus(ctx context.Context, clusterRegister *v1alpha1.ClusterRegister, opts v1.UpdateOptions) (*v1alpha1.ClusterRegister, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(clusterregistersResource, "status", c.ns, clusterRegister), &v1alpha1.ClusterRegister{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegister), err
}

// Delete takes name of the clusterRegister and deletes it. Returns an error if one occurs.
func (c *FakeClusterRegisters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(clusterregistersResource, c.ns, name), &v1alpha1.ClusterRegister{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeClusterRegisters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(clusterregistersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ClusterRegisterList{})
	return err
}

// Patch applies the patch and returns the patched clusterRegister.
func (c *FakeClusterRegisters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ClusterRegister, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(clusterregistersResource, c.ns, name, pt, data, subresources...), &v1alpha1.ClusterRegister{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ClusterRegister), err
}
