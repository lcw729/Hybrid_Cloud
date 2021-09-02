package apis

import (
	clusterregisterv1alpha1 "Hybrid_Cluster/apis/clusterRegister/v1alpha1"

	policyenginev1alpha1 "Hybrid_Cluster/apis/policy/v1alpha1"

	kubefedclusterv1alpha1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

func init() {
	// Register the types with the Scheme so the components can map objects to GroupVersionKinds and back
	AddToSchemes = append(AddToSchemes, kubefedclusterv1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, clusterregisterv1alpha1.SchemeBuilder.AddToScheme)
	AddToSchemes = append(AddToSchemes, policyenginev1alpha1.SchemeBuilder.AddToScheme)
}
