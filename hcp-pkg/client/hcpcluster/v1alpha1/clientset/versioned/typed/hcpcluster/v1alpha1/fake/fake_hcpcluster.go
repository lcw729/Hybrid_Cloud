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
	v1alpha1 "github.com/KETI-Hybrid/hcp-pkg/apis/hcpcluster/v1alpha1"
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeHCPClusters implements HCPClusterInterface
type FakeHCPClusters struct {
	Fake *FakeHcpV1alpha1
	ns   string
}

var hcpclustersResource = schema.GroupVersionResource{Group: "hcp.crd.com", Version: "v1alpha1", Resource: "hcpclusters"}

var hcpclustersKind = schema.GroupVersionKind{Group: "hcp.crd.com", Version: "v1alpha1", Kind: "HCPCluster"}

// Get takes name of the hCPCluster, and returns the corresponding hCPCluster object, and an error if there is any.
func (c *FakeHCPClusters) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.HCPCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(hcpclustersResource, c.ns, name), &v1alpha1.HCPCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HCPCluster), err
}

// List takes label and field selectors, and returns the list of HCPClusters that match those selectors.
func (c *FakeHCPClusters) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.HCPClusterList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(hcpclustersResource, hcpclustersKind, c.ns, opts), &v1alpha1.HCPClusterList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.HCPClusterList{ListMeta: obj.(*v1alpha1.HCPClusterList).ListMeta}
	for _, item := range obj.(*v1alpha1.HCPClusterList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested hCPClusters.
func (c *FakeHCPClusters) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(hcpclustersResource, c.ns, opts))

}

// Create takes the representation of a hCPCluster and creates it.  Returns the server's representation of the hCPCluster, and an error, if there is any.
func (c *FakeHCPClusters) Create(ctx context.Context, hCPCluster *v1alpha1.HCPCluster, opts v1.CreateOptions) (result *v1alpha1.HCPCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(hcpclustersResource, c.ns, hCPCluster), &v1alpha1.HCPCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HCPCluster), err
}

// Update takes the representation of a hCPCluster and updates it. Returns the server's representation of the hCPCluster, and an error, if there is any.
func (c *FakeHCPClusters) Update(ctx context.Context, hCPCluster *v1alpha1.HCPCluster, opts v1.UpdateOptions) (result *v1alpha1.HCPCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(hcpclustersResource, c.ns, hCPCluster), &v1alpha1.HCPCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HCPCluster), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeHCPClusters) UpdateStatus(ctx context.Context, hCPCluster *v1alpha1.HCPCluster, opts v1.UpdateOptions) (*v1alpha1.HCPCluster, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(hcpclustersResource, "status", c.ns, hCPCluster), &v1alpha1.HCPCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HCPCluster), err
}

// Delete takes name of the hCPCluster and deletes it. Returns an error if one occurs.
func (c *FakeHCPClusters) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(hcpclustersResource, c.ns, name), &v1alpha1.HCPCluster{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeHCPClusters) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(hcpclustersResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.HCPClusterList{})
	return err
}

// Patch applies the patch and returns the patched hCPCluster.
func (c *FakeHCPClusters) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.HCPCluster, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(hcpclustersResource, c.ns, name, pt, data, subresources...), &v1alpha1.HCPCluster{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HCPCluster), err
}
