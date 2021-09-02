package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TemplateResourceSpec defines the desired state of TemplateResource
// +k8s:openapi-gen=true
type TemplateResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Template appsv1.Deployment `json:"template" protobuf:"bytes,3,opt,name=template"`
	Replicas int32 `json:"replicas" protobuf:"varint,1,opt,name=replicas"`

        //Placement

}

// TemplateResourceStatus defines the observed state of TemplateResource
// +k8s:openapi-gen=true
type TemplateResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas int32 `json:"replicas"`
	ClusterMaps map[string]int32 `json:"clusters"`

}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplateResource is the Schema for the templateresources API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type TemplateResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TemplateResourceSpec   `json:"spec,omitempty"`
	Status TemplateResourceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TemplateResourceList contains a list of TemplateResource
type TemplateResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TemplateResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TemplateResource{}, &TemplateResourceList{})
}
