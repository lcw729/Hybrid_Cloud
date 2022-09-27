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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta2

import (
	v1beta2 "github.com/KETI-Hybrid/hcp-pkg/apis/vpa/v1beta2"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// VerticalPodAutoscalerCheckpointLister helps list VerticalPodAutoscalerCheckpoints.
// All objects returned here must be treated as read-only.
type VerticalPodAutoscalerCheckpointLister interface {
	// List lists all VerticalPodAutoscalerCheckpoints in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.VerticalPodAutoscalerCheckpoint, err error)
	// VerticalPodAutoscalerCheckpoints returns an object that can list and get VerticalPodAutoscalerCheckpoints.
	VerticalPodAutoscalerCheckpoints(namespace string) VerticalPodAutoscalerCheckpointNamespaceLister
	VerticalPodAutoscalerCheckpointListerExpansion
}

// verticalPodAutoscalerCheckpointLister implements the VerticalPodAutoscalerCheckpointLister interface.
type verticalPodAutoscalerCheckpointLister struct {
	indexer cache.Indexer
}

// NewVerticalPodAutoscalerCheckpointLister returns a new VerticalPodAutoscalerCheckpointLister.
func NewVerticalPodAutoscalerCheckpointLister(indexer cache.Indexer) VerticalPodAutoscalerCheckpointLister {
	return &verticalPodAutoscalerCheckpointLister{indexer: indexer}
}

// List lists all VerticalPodAutoscalerCheckpoints in the indexer.
func (s *verticalPodAutoscalerCheckpointLister) List(selector labels.Selector) (ret []*v1beta2.VerticalPodAutoscalerCheckpoint, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.VerticalPodAutoscalerCheckpoint))
	})
	return ret, err
}

// VerticalPodAutoscalerCheckpoints returns an object that can list and get VerticalPodAutoscalerCheckpoints.
func (s *verticalPodAutoscalerCheckpointLister) VerticalPodAutoscalerCheckpoints(namespace string) VerticalPodAutoscalerCheckpointNamespaceLister {
	return verticalPodAutoscalerCheckpointNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// VerticalPodAutoscalerCheckpointNamespaceLister helps list and get VerticalPodAutoscalerCheckpoints.
// All objects returned here must be treated as read-only.
type VerticalPodAutoscalerCheckpointNamespaceLister interface {
	// List lists all VerticalPodAutoscalerCheckpoints in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1beta2.VerticalPodAutoscalerCheckpoint, err error)
	// Get retrieves the VerticalPodAutoscalerCheckpoint from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1beta2.VerticalPodAutoscalerCheckpoint, error)
	VerticalPodAutoscalerCheckpointNamespaceListerExpansion
}

// verticalPodAutoscalerCheckpointNamespaceLister implements the VerticalPodAutoscalerCheckpointNamespaceLister
// interface.
type verticalPodAutoscalerCheckpointNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all VerticalPodAutoscalerCheckpoints in the indexer for a given namespace.
func (s verticalPodAutoscalerCheckpointNamespaceLister) List(selector labels.Selector) (ret []*v1beta2.VerticalPodAutoscalerCheckpoint, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta2.VerticalPodAutoscalerCheckpoint))
	})
	return ret, err
}

// Get retrieves the VerticalPodAutoscalerCheckpoint from the indexer for a given namespace and name.
func (s verticalPodAutoscalerCheckpointNamespaceLister) Get(name string) (*v1beta2.VerticalPodAutoscalerCheckpoint, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta2.Resource("verticalpodautoscalercheckpoint"), name)
	}
	return obj.(*v1beta2.VerticalPodAutoscalerCheckpoint), nil
}
