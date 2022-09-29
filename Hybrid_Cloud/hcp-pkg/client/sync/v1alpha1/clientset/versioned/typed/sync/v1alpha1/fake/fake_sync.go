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
	v1alpha1 "github.com/KETI-Hybrid/hcp-pkg/apis/sync/v1alpha1"
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSyncs implements SyncInterface
type FakeSyncs struct {
	Fake *FakeHcpV1alpha1
	ns   string
}

var syncsResource = schema.GroupVersionResource{Group: "hcp.crd.com", Version: "v1alpha1", Resource: "syncs"}

var syncsKind = schema.GroupVersionKind{Group: "hcp.crd.com", Version: "v1alpha1", Kind: "Sync"}

// Get takes name of the sync, and returns the corresponding sync object, and an error if there is any.
func (c *FakeSyncs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Sync, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(syncsResource, c.ns, name), &v1alpha1.Sync{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Sync), err
}

// List takes label and field selectors, and returns the list of Syncs that match those selectors.
func (c *FakeSyncs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.SyncList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(syncsResource, syncsKind, c.ns, opts), &v1alpha1.SyncList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SyncList{ListMeta: obj.(*v1alpha1.SyncList).ListMeta}
	for _, item := range obj.(*v1alpha1.SyncList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested syncs.
func (c *FakeSyncs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(syncsResource, c.ns, opts))

}

// Create takes the representation of a sync and creates it.  Returns the server's representation of the sync, and an error, if there is any.
func (c *FakeSyncs) Create(ctx context.Context, sync *v1alpha1.Sync, opts v1.CreateOptions) (result *v1alpha1.Sync, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(syncsResource, c.ns, sync), &v1alpha1.Sync{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Sync), err
}

// Update takes the representation of a sync and updates it. Returns the server's representation of the sync, and an error, if there is any.
func (c *FakeSyncs) Update(ctx context.Context, sync *v1alpha1.Sync, opts v1.UpdateOptions) (result *v1alpha1.Sync, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(syncsResource, c.ns, sync), &v1alpha1.Sync{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Sync), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSyncs) UpdateStatus(ctx context.Context, sync *v1alpha1.Sync, opts v1.UpdateOptions) (*v1alpha1.Sync, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(syncsResource, "status", c.ns, sync), &v1alpha1.Sync{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Sync), err
}

// Delete takes name of the sync and deletes it. Returns an error if one occurs.
func (c *FakeSyncs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(syncsResource, c.ns, name), &v1alpha1.Sync{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSyncs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(syncsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.SyncList{})
	return err
}

// Patch applies the patch and returns the patched sync.
func (c *FakeSyncs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Sync, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(syncsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Sync{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Sync), err
}
