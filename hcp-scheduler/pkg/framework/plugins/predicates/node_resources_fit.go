package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/framework/v1alpha1"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/hcp-scheduler/pkg/util"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type NodeResourcesFit struct{}

func (pl *NodeResourcesFit) Name() string {
	return v1alpha1.NodeResourcesFit
}

// Filter invoked at the filter extension point.
// Checks if a node has sufficient resources, such as cpu, memory, gpu, opaque int resources etc to run a pod.
// It returns a list of insufficient resources, if empty, then the node has all the resources requested by the pod.
func (p1 *NodeResourcesFit) Filter(pod *v1.Pod, status *v1alpha1.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	for _, nodeInfo := range clusterInfo.Nodes {
		insufficientResources := fitsRequest(pod, nodeInfo)
		if len(insufficientResources) == 0 {
			fmt.Println("All resources is sufficient")
			return false
		} else {
			for _, insufficientResource := range insufficientResources {
				reason := insufficientResource.Reason
				fmt.Printf("Reason: %s\n", reason)
				if reason == "Too many pods" {
					fmt.Printf("Current Pod Number: %d\n", insufficientResource.Used)
					fmt.Printf("Allowed Pod Number: %d\n", insufficientResource.Capacity)
					fmt.Println()
				} else {
					fmt.Printf("%s RequestedResources: %d\n", insufficientResource.ResourceName, insufficientResource.RequestedResources)
					fmt.Printf("%s Used: %d\n", insufficientResource.ResourceName, insufficientResource.Used)
					fmt.Printf("%s Capacity: %d\n", insufficientResource.ResourceName, insufficientResource.Capacity)
				}
			}
		}
	}
	return true
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
