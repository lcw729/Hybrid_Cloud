// +build !ignore_autogenerated

/*
Copyright 2018 The Kubernetes Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *APIResource) DeepCopyInto(out *APIResource) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new APIResource.
func (in *APIResource) DeepCopy() *APIResource {
	if in == nil {
		return nil
	}
	out := new(APIResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterCondition) DeepCopyInto(out *ClusterCondition) {
	*out = *in
	in.LastProbeTime.DeepCopyInto(&out.LastProbeTime)
	if in.LastTransitionTime != nil {
		in, out := &in.LastTransitionTime, &out.LastTransitionTime
		*out = (*in).DeepCopy()
	}
	if in.Reason != nil {
		in, out := &in.Reason, &out.Reason
		*out = new(string)
		**out = **in
	}
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterCondition.
func (in *ClusterCondition) DeepCopy() *ClusterCondition {
	if in == nil {
		return nil
	}
	out := new(ClusterCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterHealthCheckConfig) DeepCopyInto(out *ClusterHealthCheckConfig) {
	*out = *in
	if in.Period != nil {
		in, out := &in.Period, &out.Period
		*out = new(v1.Duration)
		**out = **in
	}
	if in.FailureThreshold != nil {
		in, out := &in.FailureThreshold, &out.FailureThreshold
		*out = new(int64)
		**out = **in
	}
	if in.SuccessThreshold != nil {
		in, out := &in.SuccessThreshold, &out.SuccessThreshold
		*out = new(int64)
		**out = **in
	}
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(v1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterHealthCheckConfig.
func (in *ClusterHealthCheckConfig) DeepCopy() *ClusterHealthCheckConfig {
	if in == nil {
		return nil
	}
	out := new(ClusterHealthCheckConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DurationConfig) DeepCopyInto(out *DurationConfig) {
	*out = *in
	if in.AvailableDelay != nil {
		in, out := &in.AvailableDelay, &out.AvailableDelay
		*out = new(v1.Duration)
		**out = **in
	}
	if in.UnavailableDelay != nil {
		in, out := &in.UnavailableDelay, &out.UnavailableDelay
		*out = new(v1.Duration)
		**out = **in
	}
	if in.CacheSyncTimeout != nil {
		in, out := &in.CacheSyncTimeout, &out.CacheSyncTimeout
		*out = new(v1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DurationConfig.
func (in *DurationConfig) DeepCopy() *DurationConfig {
	if in == nil {
		return nil
	}
	out := new(DurationConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FeatureGatesConfig) DeepCopyInto(out *FeatureGatesConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FeatureGatesConfig.
func (in *FeatureGatesConfig) DeepCopy() *FeatureGatesConfig {
	if in == nil {
		return nil
	}
	out := new(FeatureGatesConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FederatedTypeConfig) DeepCopyInto(out *FederatedTypeConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FederatedTypeConfig.
func (in *FederatedTypeConfig) DeepCopy() *FederatedTypeConfig {
	if in == nil {
		return nil
	}
	out := new(FederatedTypeConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FederatedTypeConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FederatedTypeConfigList) DeepCopyInto(out *FederatedTypeConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]FederatedTypeConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FederatedTypeConfigList.
func (in *FederatedTypeConfigList) DeepCopy() *FederatedTypeConfigList {
	if in == nil {
		return nil
	}
	out := new(FederatedTypeConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *FederatedTypeConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FederatedTypeConfigSpec) DeepCopyInto(out *FederatedTypeConfigSpec) {
	*out = *in
	out.TargetType = in.TargetType
	out.FederatedType = in.FederatedType
	if in.StatusType != nil {
		in, out := &in.StatusType, &out.StatusType
		*out = new(APIResource)
		**out = **in
	}
	if in.StatusCollection != nil {
		in, out := &in.StatusCollection, &out.StatusCollection
		*out = new(StatusCollectionMode)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FederatedTypeConfigSpec.
func (in *FederatedTypeConfigSpec) DeepCopy() *FederatedTypeConfigSpec {
	if in == nil {
		return nil
	}
	out := new(FederatedTypeConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FederatedTypeConfigStatus) DeepCopyInto(out *FederatedTypeConfigStatus) {
	*out = *in
	if in.StatusController != nil {
		in, out := &in.StatusController, &out.StatusController
		*out = new(ControllerStatus)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FederatedTypeConfigStatus.
func (in *FederatedTypeConfigStatus) DeepCopy() *FederatedTypeConfigStatus {
	if in == nil {
		return nil
	}
	out := new(FederatedTypeConfigStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeFedCluster) DeepCopyInto(out *KubeFedCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeFedCluster.
func (in *KubeFedCluster) DeepCopy() *KubeFedCluster {
	if in == nil {
		return nil
	}
	out := new(KubeFedCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeFedCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeFedClusterList) DeepCopyInto(out *KubeFedClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KubeFedCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeFedClusterList.
func (in *KubeFedClusterList) DeepCopy() *KubeFedClusterList {
	if in == nil {
		return nil
	}
	out := new(KubeFedClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeFedClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeFedClusterSpec) DeepCopyInto(out *KubeFedClusterSpec) {
	*out = *in
	if in.CABundle != nil {
		in, out := &in.CABundle, &out.CABundle
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	out.SecretRef = in.SecretRef
	if in.DisabledTLSValidations != nil {
		in, out := &in.DisabledTLSValidations, &out.DisabledTLSValidations
		*out = make([]TLSValidation, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeFedClusterSpec.
func (in *KubeFedClusterSpec) DeepCopy() *KubeFedClusterSpec {
	if in == nil {
		return nil
	}
	out := new(KubeFedClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeFedClusterStatus) DeepCopyInto(out *KubeFedClusterStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]ClusterCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Zones != nil {
		in, out := &in.Zones, &out.Zones
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Region != nil {
		in, out := &in.Region, &out.Region
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeFedClusterStatus.
func (in *KubeFedClusterStatus) DeepCopy() *KubeFedClusterStatus {
	if in == nil {
		return nil
	}
	out := new(KubeFedClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeFedConfig) DeepCopyInto(out *KubeFedConfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeFedConfig.
func (in *KubeFedConfig) DeepCopy() *KubeFedConfig {
	if in == nil {
		return nil
	}
	out := new(KubeFedConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeFedConfig) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeFedConfigList) DeepCopyInto(out *KubeFedConfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KubeFedConfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeFedConfigList.
func (in *KubeFedConfigList) DeepCopy() *KubeFedConfigList {
	if in == nil {
		return nil
	}
	out := new(KubeFedConfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeFedConfigList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeFedConfigSpec) DeepCopyInto(out *KubeFedConfigSpec) {
	*out = *in
	if in.ControllerDuration != nil {
		in, out := &in.ControllerDuration, &out.ControllerDuration
		*out = new(DurationConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.LeaderElect != nil {
		in, out := &in.LeaderElect, &out.LeaderElect
		*out = new(LeaderElectConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.FeatureGates != nil {
		in, out := &in.FeatureGates, &out.FeatureGates
		*out = make([]FeatureGatesConfig, len(*in))
		copy(*out, *in)
	}
	if in.ClusterHealthCheck != nil {
		in, out := &in.ClusterHealthCheck, &out.ClusterHealthCheck
		*out = new(ClusterHealthCheckConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.SyncController != nil {
		in, out := &in.SyncController, &out.SyncController
		*out = new(SyncControllerConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.StatusController != nil {
		in, out := &in.StatusController, &out.StatusController
		*out = new(StatusControllerConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeFedConfigSpec.
func (in *KubeFedConfigSpec) DeepCopy() *KubeFedConfigSpec {
	if in == nil {
		return nil
	}
	out := new(KubeFedConfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LeaderElectConfig) DeepCopyInto(out *LeaderElectConfig) {
	*out = *in
	if in.LeaseDuration != nil {
		in, out := &in.LeaseDuration, &out.LeaseDuration
		*out = new(v1.Duration)
		**out = **in
	}
	if in.RenewDeadline != nil {
		in, out := &in.RenewDeadline, &out.RenewDeadline
		*out = new(v1.Duration)
		**out = **in
	}
	if in.RetryPeriod != nil {
		in, out := &in.RetryPeriod, &out.RetryPeriod
		*out = new(v1.Duration)
		**out = **in
	}
	if in.ResourceLock != nil {
		in, out := &in.ResourceLock, &out.ResourceLock
		*out = new(ResourceLockType)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LeaderElectConfig.
func (in *LeaderElectConfig) DeepCopy() *LeaderElectConfig {
	if in == nil {
		return nil
	}
	out := new(LeaderElectConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LocalSecretReference) DeepCopyInto(out *LocalSecretReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LocalSecretReference.
func (in *LocalSecretReference) DeepCopy() *LocalSecretReference {
	if in == nil {
		return nil
	}
	out := new(LocalSecretReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatusControllerConfig) DeepCopyInto(out *StatusControllerConfig) {
	*out = *in
	if in.MaxConcurrentReconciles != nil {
		in, out := &in.MaxConcurrentReconciles, &out.MaxConcurrentReconciles
		*out = new(int64)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatusControllerConfig.
func (in *StatusControllerConfig) DeepCopy() *StatusControllerConfig {
	if in == nil {
		return nil
	}
	out := new(StatusControllerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SyncControllerConfig) DeepCopyInto(out *SyncControllerConfig) {
	*out = *in
	if in.MaxConcurrentReconciles != nil {
		in, out := &in.MaxConcurrentReconciles, &out.MaxConcurrentReconciles
		*out = new(int64)
		**out = **in
	}
	if in.AdoptResources != nil {
		in, out := &in.AdoptResources, &out.AdoptResources
		*out = new(ResourceAdoption)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SyncControllerConfig.
func (in *SyncControllerConfig) DeepCopy() *SyncControllerConfig {
	if in == nil {
		return nil
	}
	out := new(SyncControllerConfig)
	in.DeepCopyInto(out)
	return out
}
