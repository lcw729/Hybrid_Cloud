package v1alpha1

import (
	hpav2beta1 "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpav1beta2 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1beta2"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HCPHybridAutoScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPHybridAutoScalerSpec   `json:"spec,omitempty"`
	Status HCPHybridAutoScalerStatus `json:"status,omitempty"`
}

type HCPHybridAutoScalerSpec struct {
	TargetCluster  string         `json:"targetCluster"`
	WarningCount   int32          `json:"warningCount"`
	CurrentStep    string         `json:"currentStep"`
	ScalingOptions ScalingOptions `json:"scalingOptions,omitempty" protobuf:"bytes,2,opt,name=scalingoptions"`
}

type ScalingOptions struct {
	// CpaTemplate CpaTemplate                        `json:"cpaTemplate,omitempty" protobuf:"bytes,1,opt,name=cpatemplate"`
	HpaTemplate hpav2beta1.HorizontalPodAutoscaler `json:"hpaTemplate,omitempty" protobuf:"bytes,2,opt,name=hpatemplate"`
	VpaTemplate vpav1beta2.VerticalPodAutoscaler   `json:"vpaTemplate,omitempty" protobuf:"bytes,3,opt,name=hpatemplate"`
}

type HCPHybridAutoScalerStatus struct {
	ResourceStatus string                  `json:"resourceStatus"`
	LastSpec       HCPHybridAutoScalerSpec `json:"lastSpec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HCPHybridAutoScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPHybridAutoScaler `json:"items"`
}
