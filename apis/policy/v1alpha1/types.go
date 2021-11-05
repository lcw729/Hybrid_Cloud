package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HCPPolicyEngineSpec defines the desired state of HCPPolicyEngine
// +k8s:openapi-gen=true

type HCPPolicyTemplate struct {
	Spec HCPPolicyTemplateSpec `json:"spec"`
}

type HCPPolicyTemplateSpec struct {
	TargetController HCPPolicyTartgetController `json:"targetController"`
	Policies         []HCPPolicies              `json:"policies"`
}

type HCPPolicyTartgetController struct {
	Kind string `json:"kind"`
}

type HCPPolicies struct {
	Type  string   `json:"type"`
	Value []string `json:"value"`
}

type HCPPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	//Template - 생성
	Template HCPPolicyTemplate `json:"template"`
	/*Template struct {
		Spec struct {
			TargetController struct {
				Kind string `json:"kind"`
			} `json:"targetController"`
			Policies []struct {
				Type string `json:"type"`
				Value string `json:"value"`
			} `json:"policies"`
		} `json:"spec"`
	} `json:"template"`*/
	RangeOfApplication string `json:"rangeOfApplication"`
	PolicyStatus       string `json:"policyStatus"`
	//Placement

}

// HCPPolicyStatus defines the observed state of HCPPolicy
// +k8s:openapi-gen=true
type HCPPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Replicas    int32            `json:"replicas"`
	ClusterMaps map[string]int32 `json:"clusters"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPPolicy is the Schema for the HCPpolicys API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type HCPPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HCPPolicySpec   `json:"spec,omitempty"`
	Status HCPPolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HCPPolicyList contains a list of HCPPolicy
type HCPPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HCPPolicy `json:"items"`
}
