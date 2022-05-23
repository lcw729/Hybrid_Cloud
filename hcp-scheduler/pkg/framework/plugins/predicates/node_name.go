package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/framework/v1alpha1"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type NodeName struct{}

func (pl *NodeName) Name() string {
	return v1alpha1.NodeName
}

// Filter invoked at the filter extension point.
func (pl *NodeName) Filter(pod *v1.Pod, status *v1alpha1.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	for _, nodeInfo := range clusterInfo.Nodes {
		if nodeInfo.Node == nil {
			fmt.Println("node not found")
			//nodeInfo.FilterNode()
			//clusterInfo.MinusOneAvailableNodes()
			continue
		}
		if !Fits(pod, nodeInfo) {
			fmt.Println("Node Name is unmatched")
			//nodeInfo.FilterNode()
			//clusterInfo.MinusOneAvailableNodes()
		} else {
			return false
		}
	}
	return true
}

// Fits actually checks if the pod fits the node.
func Fits(pod *v1.Pod, nodeInfo *resourceinfo.NodeInfo) bool {
	return len(pod.Spec.NodeName) == 0 || pod.Spec.NodeName == nodeInfo.Node.Name
}
