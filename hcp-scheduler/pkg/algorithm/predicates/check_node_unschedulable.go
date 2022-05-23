package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/priorities"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

const (
	// CheckNodeUnschedulablePred defines the name of predicate CheckNodeUnschedulablePredicate.
	CheckNodeUnschedulablePred = "CheckNodeUnschedulable"
)

// CheckNodeUnschedulablePredicate checks if a pod can be scheduled on a node with Unschedulable spec.
func CheckNodeUnschedulable(pod *v1.Pod, clusterInfo *resourceinfo.ClusterInfo) {
	var temp []*resourceinfo.NodeInfo
	for _, nodeInfo := range clusterInfo.Nodes {
		fmt.Println("=============================")
		fmt.Println("<<", nodeInfo.NodeName, ">>")
		// If pod tolerate unschedulable taint, it's also tolerate `node.Spec.Unschedulable`.
		podToleratesUnschedulable := priorities.TolerationsTolerateTaint(pod.Spec.Tolerations, &v1.Taint{
			Key:    v1.TaintNodeUnschedulable,
			Effect: v1.TaintEffectNoSchedule,
		})

		// TODO (k82cn): deprecates `node.Spec.Unschedulable` in 1.13.
		if nodeInfo.Node.Spec.Unschedulable && !podToleratesUnschedulable {
			fmt.Println("")
			temp = append(temp, nodeInfo)
		} else {
			fmt.Println("Node has Unschedulable Taint.\nbut, Pod hasn't Unschedulable Toleration.")
		}
	}
	clusterInfo.Nodes = temp
}
