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
	v1beta2 "hcp-pkg/apis/vpa/v1beta2"
	"context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVerticalPodAutoscalerCheckpoints implements VerticalPodAutoscalerCheckpointInterface
type FakeVerticalPodAutoscalerCheckpoints struct {
	Fake *FakeAutoscalingV1beta2
	ns   string
}

var verticalpodautoscalercheckpointsResource = schema.GroupVersionResource{Group: "autoscaling.k8s.io", Version: "v1beta2", Resource: "verticalpodautoscalercheckpoints"}

var verticalpodautoscalercheckpointsKind = schema.GroupVersionKind{Group: "autoscaling.k8s.io", Version: "v1beta2", Kind: "VerticalPodAutoscalerCheckpoint"}

// Get takes name of the verticalPodAutoscalerCheckpoint, and returns the corresponding verticalPodAutoscalerCheckpoint object, and an error if there is any.
func (c *FakeVerticalPodAutoscalerCheckpoints) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta2.VerticalPodAutoscalerCheckpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(verticalpodautoscalercheckpointsResource, c.ns, name), &v1beta2.VerticalPodAutoscalerCheckpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.VerticalPodAutoscalerCheckpoint), err
}

// List takes label and field selectors, and returns the list of VerticalPodAutoscalerCheckpoints that match those selectors.
func (c *FakeVerticalPodAutoscalerCheckpoints) List(ctx context.Context, opts v1.ListOptions) (result *v1beta2.VerticalPodAutoscalerCheckpointList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(verticalpodautoscalercheckpointsResource, verticalpodautoscalercheckpointsKind, c.ns, opts), &v1beta2.VerticalPodAutoscalerCheckpointList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta2.VerticalPodAutoscalerCheckpointList{ListMeta: obj.(*v1beta2.VerticalPodAutoscalerCheckpointList).ListMeta}
	for _, item := range obj.(*v1beta2.VerticalPodAutoscalerCheckpointList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested verticalPodAutoscalerCheckpoints.
func (c *FakeVerticalPodAutoscalerCheckpoints) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(verticalpodautoscalercheckpointsResource, c.ns, opts))

}

// Create takes the representation of a verticalPodAutoscalerCheckpoint and creates it.  Returns the server's representation of the verticalPodAutoscalerCheckpoint, and an error, if there is any.
func (c *FakeVerticalPodAutoscalerCheckpoints) Create(ctx context.Context, verticalPodAutoscalerCheckpoint *v1beta2.VerticalPodAutoscalerCheckpoint, opts v1.CreateOptions) (result *v1beta2.VerticalPodAutoscalerCheckpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(verticalpodautoscalercheckpointsResource, c.ns, verticalPodAutoscalerCheckpoint), &v1beta2.VerticalPodAutoscalerCheckpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.VerticalPodAutoscalerCheckpoint), err
}

// Update takes the representation of a verticalPodAutoscalerCheckpoint and updates it. Returns the server's representation of the verticalPodAutoscalerCheckpoint, and an error, if there is any.
func (c *FakeVerticalPodAutoscalerCheckpoints) Update(ctx context.Context, verticalPodAutoscalerCheckpoint *v1beta2.VerticalPodAutoscalerCheckpoint, opts v1.UpdateOptions) (result *v1beta2.VerticalPodAutoscalerCheckpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(verticalpodautoscalercheckpointsResource, c.ns, verticalPodAutoscalerCheckpoint), &v1beta2.VerticalPodAutoscalerCheckpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.VerticalPodAutoscalerCheckpoint), err
}

// Delete takes name of the verticalPodAutoscalerCheckpoint and deletes it. Returns an error if one occurs.
func (c *FakeVerticalPodAutoscalerCheckpoints) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(verticalpodautoscalercheckpointsResource, c.ns, name), &v1beta2.VerticalPodAutoscalerCheckpoint{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVerticalPodAutoscalerCheckpoints) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(verticalpodautoscalercheckpointsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta2.VerticalPodAutoscalerCheckpointList{})
	return err
}

// Patch applies the patch and returns the patched verticalPodAutoscalerCheckpoint.
func (c *FakeVerticalPodAutoscalerCheckpoints) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta2.VerticalPodAutoscalerCheckpoint, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(verticalpodautoscalercheckpointsResource, c.ns, name, pt, data, subresources...), &v1beta2.VerticalPodAutoscalerCheckpoint{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta2.VerticalPodAutoscalerCheckpoint), err
}
