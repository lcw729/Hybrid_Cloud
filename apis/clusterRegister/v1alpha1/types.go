package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterRegisterSpec struct {
	Platform      string `json:"platform" protobuf:"bytes,1,opt,name=platform"`
	Resourcegroup string `json:"resourcegroup" protobuf:"bytes,2,opt,name=resourcegroup"`
	Projectid     string `json:"projectid" protobuf:"bytes,3,opt,name=projectid"`
	Clustername   string `json:"clustername" protobuf:"bytes,4,opt,name=clustername"`
	Region        string `json:"region" protobuf:"bytes,5,opt,name=region"`
}

type ClusterRegister struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// changed
	Spec ClusterRegisterSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

type ClusterRegisterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ClusterRegister `json:"items"`
}
