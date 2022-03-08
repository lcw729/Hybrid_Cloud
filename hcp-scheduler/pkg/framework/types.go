package framework

import (
	"Hybrid_Cloud/pkg/apis/hcpcluster/v1alpha1"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
)

// ActionType is an integer to represent one type of resource change.
// Different ActionTypes can be bit-wised to compose new semantics.
type ActionType int64

// Constants for ActionTypes.
const (
	Add    ActionType = 1 << iota // 1
	Delete                        // 10
	// UpdateNodeXYZ is only applicable for Node events.
	UpdateNodeAllocatable // 100
	UpdateNodeLabel       // 1000
	UpdateNodeTaint       // 10000
	UpdateNodeCondition   // 100000

	All ActionType = 1<<iota - 1 // 111111

	// Use the general Update type if you don't either know or care the specific sub-Update type to use.
	Update = UpdateNodeAllocatable | UpdateNodeLabel | UpdateNodeTaint | UpdateNodeCondition
)

// GVK is short for group/version/kind, which can uniquely represent a particular API resource.
type GVK string

// Constants for GVKs.
const (
	Pod                   GVK = "Pod"
	Node                  GVK = "Node"
	PersistentVolume      GVK = "PersistentVolume"
	PersistentVolumeClaim GVK = "PersistentVolumeClaim"
	Service               GVK = "Service"
	StorageClass          GVK = "storage.k8s.io/StorageClass"
	CSINode               GVK = "storage.k8s.io/CSINode"
	CSIDriver             GVK = "storage.k8s.io/CSIDriver"
	CSIStorageCapacity    GVK = "storage.k8s.io/CSIStorageCapacity"
	WildCard              GVK = "*"
)

// ClusterEvent abstracts how a system resource's state gets changed.
// Resource represents the standard API resources such as Pod, Node, etc.
// ActionType denotes the specific change such as Add, Update or Delete.
type ClusterEvent struct {
	Resource   GVK
	ActionType ActionType
	Label      string
}

type ClusterInfo struct {
	ClusterName        string
	Nodes              []*NodeInfo
	RequestedResources *Resources
	AllocableResources *Resources
	CapacityResources  *Resources
	ClusterList        *v1alpha1.HCPClusterList
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

// NodeInfo is node level aggregated information.
type NodeInfo struct {
	ClusterName string
	NodeName    string
	// Overall node information.
	node *v1.Node

	// Pods running on the node.
	Pods []*PodInfo

	// The subset of pods with affinity.
	PodsWithAffinity []*PodInfo

	// The subset of pods with required anti-affinity.
	PodsWithRequiredAntiAffinity []*PodInfo

	// Total requested resources of all pods on this node.
	RequestedResources   *Resource
	AllocatableResources *Resource
	CapacityResources    *Resource
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

// IsWildCard returns true if ClusterEvent follows WildCard semantics
func (ce ClusterEvent) IsWildCard() bool {
	return ce.Resource == WildCard && ce.ActionType == All
}

// QueuedPodInfo is a Pod wrapper with additional information related to
// the pod's status in the scheduling queue, such as the timestamp when
// it's added to the queue.
type QueuedPodInfo struct {
	*PodInfo
	// The time pod added to the scheduling queue.
	Timestamp time.Time
	// Number of schedule attempts before successfully scheduled.
	// It's used to record the # attempts metric.
	Attempts int
	// The time when the pod is added to the queue for the first time. The pod may be added
	// back to the queue multiple times before it's successfully scheduled.
	// It shouldn't be updated once initialized. It's used to record the e2e scheduling
	// latency for a pod.
	InitialAttemptTimestamp time.Time
	// If a Pod failed in a scheduling cycle, record the plugin names it failed by.
	UnschedulablePlugins sets.String
}

// DeepCopy returns a deep copy of the QueuedPodInfo object.
func (pqi *QueuedPodInfo) DeepCopy() *QueuedPodInfo {
	return &QueuedPodInfo{
		PodInfo:                 pqi.PodInfo.DeepCopy(),
		Timestamp:               pqi.Timestamp,
		Attempts:                pqi.Attempts,
		InitialAttemptTimestamp: pqi.InitialAttemptTimestamp,
	}
}

// DeepCopy returns a deep copy of the PodInfo object.
func (pi *PodInfo) DeepCopy() *PodInfo {
	return &PodInfo{
		Pod:                        pi.Pod.DeepCopy(),
		RequiredAffinityTerms:      pi.RequiredAffinityTerms,
		RequiredAntiAffinityTerms:  pi.RequiredAntiAffinityTerms,
		PreferredAffinityTerms:     pi.PreferredAffinityTerms,
		PreferredAntiAffinityTerms: pi.PreferredAntiAffinityTerms,
		ParseError:                 pi.ParseError,
	}
}

// Update creates a full new PodInfo by default. And only updates the pod when the PodInfo
// has been instantiated and the passed pod is the exact same one as the original pod.
func (pi *PodInfo) Update(pod *v1.Pod) {
	if pod != nil && pi.Pod != nil && pi.Pod.UID == pod.UID {
		// PodInfo includes immutable information, and so it is safe to update the pod in place if it is
		// the exact same pod
		pi.Pod = pod
		return
	}
	var preferredAffinityTerms []v1.WeightedPodAffinityTerm
	var preferredAntiAffinityTerms []v1.WeightedPodAffinityTerm
	if affinity := pod.Spec.Affinity; affinity != nil {
		if a := affinity.PodAffinity; a != nil {
			preferredAffinityTerms = a.PreferredDuringSchedulingIgnoredDuringExecution
		}
		if a := affinity.PodAntiAffinity; a != nil {
			preferredAntiAffinityTerms = a.PreferredDuringSchedulingIgnoredDuringExecution
		}
	}

	// Attempt to parse the affinity terms
	var parseErrs []error
	requiredAffinityTerms, err := getAffinityTerms(pod, getPodAffinityTerms(pod.Spec.Affinity))
	if err != nil {
		parseErrs = append(parseErrs, fmt.Errorf("requiredAffinityTerms: %w", err))
	}
	requiredAntiAffinityTerms, err := getAffinityTerms(pod,
		getPodAntiAffinityTerms(pod.Spec.Affinity))
	if err != nil {
		parseErrs = append(parseErrs, fmt.Errorf("requiredAntiAffinityTerms: %w", err))
	}
	weightedAffinityTerms, err := getWeightedAffinityTerms(pod, preferredAffinityTerms)
	if err != nil {
		parseErrs = append(parseErrs, fmt.Errorf("preferredAffinityTerms: %w", err))
	}
	weightedAntiAffinityTerms, err := getWeightedAffinityTerms(pod, preferredAntiAffinityTerms)
	if err != nil {
		parseErrs = append(parseErrs, fmt.Errorf("preferredAntiAffinityTerms: %w", err))
	}

	pi.Pod = pod
	pi.RequiredAffinityTerms = requiredAffinityTerms
	pi.RequiredAntiAffinityTerms = requiredAntiAffinityTerms
	pi.PreferredAffinityTerms = weightedAffinityTerms
	pi.PreferredAntiAffinityTerms = weightedAntiAffinityTerms
	pi.ParseError = utilerrors.NewAggregate(parseErrs)
}

// Matches returns true if the pod matches the label selector and namespaces or namespace selector.
func (at *AffinityTerm) Matches(pod *v1.Pod, nsLabels labels.Set) bool {
	if at.Namespaces.Has(pod.Namespace) || at.NamespaceSelector.Matches(nsLabels) {
		return at.Selector.Matches(labels.Set(pod.Labels))
	}
	return false
}

// Diagnosis records the details to diagnose a scheduling failure.
type Diagnosis struct {
	NodeToStatusMap      NodeToStatusMap
	UnschedulablePlugins sets.String
	// PostFilterMsg records the messages returned from PostFilterPlugins.
	PostFilterMsg string
}

// FitError describes a fit error of a pod.
type FitError struct {
	Pod         *v1.Pod
	NumAllNodes int
	Diagnosis   Diagnosis
}

const (
	// NoNodeAvailableMsg is used to format message when no nodes available.
	NoNodeAvailableMsg = "0/%v nodes are available"
)

// Error returns detailed information of why the pod failed to fit on each node
func (f *FitError) Error() string {
	reasons := make(map[string]int)
	for _, status := range f.Diagnosis.NodeToStatusMap {
		for _, reason := range status.Reasons() {
			reasons[reason]++
		}
	}

	sortReasonsHistogram := func() []string {
		var reasonStrings []string
		for k, v := range reasons {
			reasonStrings = append(reasonStrings, fmt.Sprintf("%v %v", v, k))
		}
		sort.Strings(reasonStrings)
		return reasonStrings
	}
	reasonMsg := fmt.Sprintf(NoNodeAvailableMsg+": %v.", f.NumAllNodes, strings.Join(sortReasonsHistogram(), ", "))
	postFilterMsg := f.Diagnosis.PostFilterMsg
	if postFilterMsg != "" {
		reasonMsg += " " + postFilterMsg
	}
	return reasonMsg
}

func newAffinityTerm(pod *v1.Pod, term *v1.PodAffinityTerm) (*AffinityTerm, error) {
	selector, err := metav1.LabelSelectorAsSelector(term.LabelSelector)
	if err != nil {
		return nil, err
	}

	namespaces := getNamespacesFromPodAffinityTerm(pod, term)
	nsSelector, err := metav1.LabelSelectorAsSelector(term.NamespaceSelector)
	if err != nil {
		return nil, err
	}

	return &AffinityTerm{Namespaces: namespaces, Selector: selector, TopologyKey: term.TopologyKey, NamespaceSelector: nsSelector}, nil
}

// getAffinityTerms receives a Pod and affinity terms and returns the namespaces and
// selectors of the terms.
func getAffinityTerms(pod *v1.Pod, v1Terms []v1.PodAffinityTerm) ([]AffinityTerm, error) {
	if v1Terms == nil {
		return nil, nil
	}

	var terms []AffinityTerm
	for i := range v1Terms {
		t, err := newAffinityTerm(pod, &v1Terms[i])
		if err != nil {
			// We get here if the label selector failed to process
			return nil, err
		}
		terms = append(terms, *t)
	}
	return terms, nil
}

// getWeightedAffinityTerms returns the list of processed affinity terms.
func getWeightedAffinityTerms(pod *v1.Pod, v1Terms []v1.WeightedPodAffinityTerm) ([]WeightedAffinityTerm, error) {
	if v1Terms == nil {
		return nil, nil
	}

	var terms []WeightedAffinityTerm
	for i := range v1Terms {
		t, err := newAffinityTerm(pod, &v1Terms[i].PodAffinityTerm)
		if err != nil {
			// We get here if the label selector failed to process
			return nil, err
		}
		terms = append(terms, WeightedAffinityTerm{AffinityTerm: *t, Weight: v1Terms[i].Weight})
	}
	return terms, nil
}

// NewPodInfo returns a new PodInfo.
func NewPodInfo(pod *v1.Pod) *PodInfo {
	pInfo := &PodInfo{}
	pInfo.Update(pod)
	return pInfo
}

func getPodAffinityTerms(affinity *v1.Affinity) (terms []v1.PodAffinityTerm) {
	if affinity != nil && affinity.PodAffinity != nil {
		if len(affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution) != 0 {
			terms = affinity.PodAffinity.RequiredDuringSchedulingIgnoredDuringExecution
		}
		// TODO: Uncomment this block when implement RequiredDuringSchedulingRequiredDuringExecution.
		//if len(affinity.PodAffinity.RequiredDuringSchedulingRequiredDuringExecution) != 0 {
		//	terms = append(terms, affinity.PodAffinity.RequiredDuringSchedulingRequiredDuringExecution...)
		//}
	}
	return terms
}

func getPodAntiAffinityTerms(affinity *v1.Affinity) (terms []v1.PodAffinityTerm) {
	if affinity != nil && affinity.PodAntiAffinity != nil {
		if len(affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution) != 0 {
			terms = affinity.PodAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution
		}
		// TODO: Uncomment this block when implement RequiredDuringSchedulingRequiredDuringExecution.
		//if len(affinity.PodAntiAffinity.RequiredDuringSchedulingRequiredDuringExecution) != 0 {
		//	terms = append(terms, affinity.PodAntiAffinity.RequiredDuringSchedulingRequiredDuringExecution...)
		//}
	}
	return terms
}

// returns a set of names according to the namespaces indicated in podAffinityTerm.
// If namespaces is empty it considers the given pod's namespace.
func getNamespacesFromPodAffinityTerm(pod *v1.Pod, podAffinityTerm *v1.PodAffinityTerm) sets.String {
	names := sets.String{}
	if len(podAffinityTerm.Namespaces) == 0 && podAffinityTerm.NamespaceSelector == nil {
		names.Insert(pod.Namespace)
	} else {
		names.Insert(podAffinityTerm.Namespaces...)
	}
	return names
}

// ImageStateSummary provides summarized information about the state of an image.
type ImageStateSummary struct {
	// Size of the image
	Size int64
	// Used to track how many nodes have this image
	NumNodes int
}

// Clone returns a copy of this resource.
func (r *Resource) Clone() *Resource {
	res := &Resource{
		MilliCPU:         r.MilliCPU,
		Memory:           r.Memory,
		AllowedPodNumber: r.AllowedPodNumber,
		EphemeralStorage: r.EphemeralStorage,
	}
	if r.ScalarResources != nil {
		res.ScalarResources = make(map[v1.ResourceName]int64)
		for k, v := range r.ScalarResources {
			res.ScalarResources[k] = v
		}
	}
	return res
}

// AddScalar adds a resource by a scalar value of this resource.
func (r *Resource) AddScalar(name v1.ResourceName, quantity int64) {
	r.SetScalar(name, r.ScalarResources[name]+quantity)
}

// SetScalar sets a resource by a scalar value of this resource.
func (r *Resource) SetScalar(name v1.ResourceName, quantity int64) {
	// Lazily allocate scalar resource map.
	if r.ScalarResources == nil {
		r.ScalarResources = map[v1.ResourceName]int64{}
	}
	r.ScalarResources[name] = quantity
}

// Node returns overall information about this node.
func (n *NodeInfo) Node() *v1.Node {
	if n == nil {
		return nil
	}
	return n.node
}

// FilterOutPods receives a list of pods and filters out those whose node names
// are equal to the node of this NodeInfo, but are not found in the pods of this NodeInfo.
//
// Preemption logic simulates removal of pods on a node by removing them from the
// corresponding NodeInfo. In order for the simulation to work, we call this method
// on the pods returned from SchedulerCache, so that predicate functions see
// only the pods that are not removed from the NodeInfo.
func (n *NodeInfo) FilterOutPods(pods []*v1.Pod) []*v1.Pod {
	node := n.Node()
	if node == nil {
		return pods
	}
	filtered := make([]*v1.Pod, 0, len(pods))
	for _, p := range pods {
		if p.Spec.NodeName != node.Name {
			filtered = append(filtered, p)
			continue
		}
		// If pod is on the given node, add it to 'filtered' only if it is present in nodeInfo.
		podKey, err := GetPodKey(p)
		if err != nil {
			continue
		}
		for _, np := range n.Pods {
			npodkey, _ := GetPodKey(np.Pod)
			if npodkey == podKey {
				filtered = append(filtered, p)
				break
			}
		}
	}
	return filtered
}

// GetPodKey returns the string key of a pod.
func GetPodKey(pod *v1.Pod) (string, error) {
	uid := string(pod.UID)
	if len(uid) == 0 {
		return "", errors.New("cannot get cache key for pod with empty UID")
	}
	return uid, nil
}
