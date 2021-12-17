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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	hcppolicyv1alpha1 "Hybrid_Cluster/pkg/apis/hcppolicy/v1alpha1"
	versioned "Hybrid_Cluster/pkg/policy/v1alpha1/clientset/versioned"
	internalinterfaces "Hybrid_Cluster/pkg/policy/v1alpha1/informers/externalversions/internalinterfaces"
	v1alpha1 "Hybrid_Cluster/pkg/policy/v1alpha1/listers/hcppolicy/v1alpha1"
	"context"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// HCPPolicyInformer provides access to a shared informer and lister for
// HCPPolicies.
type HCPPolicyInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.HCPPolicyLister
}

type hCPPolicyInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewHCPPolicyInformer constructs a new informer for HCPPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewHCPPolicyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredHCPPolicyInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredHCPPolicyInformer constructs a new informer for HCPPolicy type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredHCPPolicyInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HcpV1alpha1().HCPPolicies(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HcpV1alpha1().HCPPolicies(namespace).Watch(context.TODO(), options)
			},
		},
		&hcppolicyv1alpha1.HCPPolicy{},
		resyncPeriod,
		indexers,
	)
}

func (f *hCPPolicyInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredHCPPolicyInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *hCPPolicyInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&hcppolicyv1alpha1.HCPPolicy{}, f.defaultInformer)
}

func (f *hCPPolicyInformer) Lister() v1alpha1.HCPPolicyLister {
	return v1alpha1.NewHCPPolicyLister(f.Informer().GetIndexer())
}
