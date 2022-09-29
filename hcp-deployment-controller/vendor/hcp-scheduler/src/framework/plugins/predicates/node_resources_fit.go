package predicates

import (
	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/resourceinfo"
	"hcp-scheduler/src/util"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type NodeResourcesFit struct{}

func (pl *NodeResourcesFit) Name() string {
	return plugins.NodeResourcesFit
}

// Filter invoked at the filter extension point.
// Checks if a node has sufficient resources, such as cpu, memory, gpu, opaque int resources etc to run a pod.
// It returns a list of insufficient resources, if empty, then the node has all the resources requested by the pod.
func (p1 *NodeResourcesFit) Filter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	for _, nodeInfo := range clusterInfo.Nodes {
		insufficientResources := fitsRequest(pod, nodeInfo)
		if len(insufficientResources) == 0 {
			klog.Infoln("All resources is sufficient")
			return false
		} else {
			for _, insufficientResource := range insufficientResources {
				reason := insufficientResource.Reason
				klog.Infof("Reason: %s\n", reason)
				if reason == "Too many pods" {
					klog.Infof("Current Pod Number: %d\n", insufficientResource.Used)
					klog.Infof("Allowed Pod Number: %d\n", insufficientResource.Capacity)
					klog.Infoln()
				} else {
					klog.Infof("%s RequestedResources: %d\n", insufficientResource.ResourceName, insufficientResource.RequestedResources)
					klog.Infof("%s Used: %d\n", insufficientResource.ResourceName, insufficientResource.Used)
					klog.Infof("%s Capacity: %d\n", insufficientResource.ResourceName, insufficientResource.Capacity)
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
