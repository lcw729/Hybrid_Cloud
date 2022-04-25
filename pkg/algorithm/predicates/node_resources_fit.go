package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/hcp-scheduler/pkg/util"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// Filter invoked at the filter extension point.
// Checks if a node has sufficient resources, such as cpu, memory, gpu, opaque int resources etc to run a pod.
// It returns a list of insufficient resources, if empty, then the node has all the resources requested by the pod.
func NodeResourcesFit(pod *v1.Pod, clusterInfo *resourceinfo.ClusterInfo) {
	var temp []*resourceinfo.NodeInfo
	for _, nodeInfo := range clusterInfo.Nodes {
		fmt.Println(nodeInfo.NodeName)
		insufficientResources := fitsRequest(pod, nodeInfo)
		fmt.Println(insufficientResources)
		if insufficientResources == nil {
			temp = append(temp, nodeInfo)
		}
	}
	clusterInfo.Nodes = temp
}

// InsufficientResource describes what kind of resource limit is hit and caused the pod to not fit the node.
type InsufficientResource struct {
	ResourceName v1.ResourceName
	// We explicitly have a parameter for reason to avoid formatting a message on the fly
	// for common resources, which is expensive for cluster autoscaler simulations.
	Reason             string
	RequestedResources int64
	Used               int64
	Capacity           int64
}

func fitsRequest(pod *v1.Pod, nodeInfo *resourceinfo.NodeInfo) []InsufficientResource {
	insufficientResources := make([]InsufficientResource, 0, 4)

	allowedPodNumber := nodeInfo.AllocatableResources.AllowedPodNumber
	fmt.Println(allowedPodNumber)
	if len(nodeInfo.Pods)+1 > allowedPodNumber {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourcePods,
			"Too many pods",
			1,
			int64(len(nodeInfo.Pods)),
			int64(allowedPodNumber),
		})
	}

	podRequest := util.CreateResourceToValueMapPO(pod)
	if podRequest[v1.ResourceCPU] == 0 &&
		podRequest[v1.ResourceMemory] == 0 &&
		podRequest[v1.ResourceEphemeralStorage] == 0 {
		return insufficientResources
	}

	fmt.Println("here")
	if podRequest[v1.ResourceCPU] > (nodeInfo.AllocatableResources.MilliCPU - nodeInfo.RequestedResources.MilliCPU) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceCPU,
			"Insufficient cpu",
			podRequest[v1.ResourceCPU],
			nodeInfo.RequestedResources.MilliCPU,
			nodeInfo.AllocatableResources.MilliCPU,
		})
	}
	if podRequest[v1.ResourceMemory] > (nodeInfo.AllocatableResources.Memory - nodeInfo.RequestedResources.Memory) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceMemory,
			"Insufficient memory",
			podRequest[v1.ResourceMemory],
			nodeInfo.RequestedResources.Memory,
			nodeInfo.AllocatableResources.Memory,
		})
	}
	if podRequest[v1.ResourceEphemeralStorage] > (nodeInfo.AllocatableResources.EphemeralStorage - nodeInfo.RequestedResources.EphemeralStorage) {
		insufficientResources = append(insufficientResources, InsufficientResource{
			v1.ResourceEphemeralStorage,
			"Insufficient ephemeral-storage",
			podRequest[v1.ResourceEphemeralStorage],
			nodeInfo.RequestedResources.EphemeralStorage,
			nodeInfo.AllocatableResources.EphemeralStorage,
		})
	}

	return insufficientResources
}
