package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HCPCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPClusterSpec   `json:"spec,omitempty"`
	Status HCPClusterStatus `json:"status,omitempty"`
}

type HCPClusterSpec struct {
	Platform      string `json:"platform"`
	Region        string `json:"region"`
	Name          string `json:"name"`
	Resourcegroup string `json:"resourcegroup"`
	ProjectId     string `json:"projectid"`
}

type HCPClusterStatus struct {
	Join bool `json:"join"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HCPClusterList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Items []HCPCluster `json:"items"`
}
