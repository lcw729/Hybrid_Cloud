package resourceinfo

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
)

type ClusterInfoList []*ClusterInfo

type ClusterInfo struct {
	ClusterName  string
	IsFiltered   bool
	ClusterScore int32
	// AvailableNodes int
	Nodes              []*NodeInfo
	RequestedResources *Resources
	AllocableResources *Resources
	CapacityResources  *Resources
}

/*
func (c *ClusterInfo) MinusOneAvailableNodes() {
	c.AvailableNodes--
}

func (c *ClusterInfo) IsAnyNodes() bool {
	if c.AvailableNodes == 0 {
		return false
	} else {
		return true
	}
}
*/

// ProtocolPort represents a protocol port pair, e.g. tcp:80.
type ProtocolPort struct {
	Protocol string
	Port     int32
}

// HostPortInfo stores mapping from ip to a set of ProtocolPort
type HostPortInfo map[string]map[ProtocolPort]struct{}

// NodeInfo is node level aggregated information.
type NodeInfo struct {
	ClusterName string
	NodeName    string
	// Overall node information.
	Node *v1.Node
	// Ports allocated on the node.
	UsedPorts HostPortInfo

	// Pods running on the node.
	Pods []*PodInfo

	// The subset of pods with affinity.
	PodsWithAffinity []*PodInfo

	// The subset of pods with required anti-affinity.
	PodsWithRequiredAntiAffinity []*PodInfo

	ImageStates map[string]*ImageStateSummary

	// Total requested resources of all pods on this node.
	RequestedResources   Resource
	AllocatableResources *Resource
	CapacityResources    *Resource
}

// PodInfo is a wrapper to a Pod with additional pre-computed information to
// accelerate processing. This information is typically immutable (e.g., pre-processed
// inter-pod affinity selectors).
type PodInfo struct {
	ClusterName string
	NodeName    string
	PodName     string

	Pod *v1.Pod

	RequestedResources *Resources
	AllocableResources *Resources
	CapacityResources  *Resources

	RequiredAffinityTerms      []AffinityTerm
	RequiredAntiAffinityTerms  []AffinityTerm
	PreferredAffinityTerms     []WeightedAffinityTerm
	PreferredAntiAffinityTerms []WeightedAffinityTerm
	ParseError                 error
}

// AffinityTerm is a processed version of v1.PodAffinityTerm.
type AffinityTerm struct {
	Namespaces        sets.String
	Selector          labels.Selector
	TopologyKey       string
	NamespaceSelector labels.Selector
}

// WeightedAffinityTerm is a "processed" representation of v1.WeightedAffinityTerm.
type WeightedAffinityTerm struct {
	AffinityTerm
	Weight int32
}

// ImageStateSummary provides summarized information about the state of an image.
type ImageStateSummary struct {
	// Size of the image
	Size int64
	// Used to track how many nodes have this image
	NumNodes int
}

// Resource is a collection of compute resource.
type Resource struct {
	MilliCPU         int64
	Memory           int64
	EphemeralStorage int64
	// We store allowedPodNumber (which is Node.Status.Allocatable.Pods().Value())
	// explicitly as int, to avoid conversions and improve performance.
	AllowedPodNumber int
	// ScalarResources
	ScalarResources map[v1.ResourceName]int64
}

type Resources struct {
	CPU     string
	Memory  string
	Fs      string
	Network string
}
