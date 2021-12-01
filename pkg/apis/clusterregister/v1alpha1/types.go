package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterRegister struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterRegisterSpec   `json:"spec,omitempty"`
	Status ClusterRegisterStatus `json:"status,omitempty"`
}

type ClusterRegisterSpec struct {
	Platform      string `json:"platform"`
	Region        string `json:"region"`
	Name          string `json:"name"`
	Resourcegroup string `json:"resourcegroup"`
	ProjectId     string `json:"projectid"`
}

type ClusterRegisterStatus struct {
	Join bool `json:"join"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ClusterRegisterList struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Items []ClusterRegister `json:"items"`
}
