package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

func CheckMatchNodeName(pod *v1.Pod, clusterInfo *resourceinfo.ClusterInfo) {
	var temp []*resourceinfo.NodeInfo
	for _, nodeInfo := range clusterInfo.Nodes {
		fmt.Println("=============================")
		fmt.Println("<<", nodeInfo.NodeName, ">>")

		if MatchNodeName(pod, nodeInfo) {
			temp = append(temp, nodeInfo)
		} else {
			fmt.Println("Node Name is unmatched")
		}
	}
	clusterInfo.Nodes = temp
}

// Filter invoked at the filter extension point.
func MatchNodeName(pod *v1.Pod, nodeInfo *resourceinfo.NodeInfo) bool {
	if nodeInfo.Node == nil {
		return false
	}
	if !Fits(pod, nodeInfo) {
		return false
	}
	return true
}

// Fits actually checks if the pod fits the node.
func Fits(pod *v1.Pod, nodeInfo *resourceinfo.NodeInfo) bool {
	return len(pod.Spec.NodeName) == 0 || pod.Spec.NodeName == nodeInfo.Node.Name
}
