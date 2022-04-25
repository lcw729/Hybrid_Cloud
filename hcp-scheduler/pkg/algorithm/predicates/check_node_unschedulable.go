package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/priorities"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"

	v1 "k8s.io/api/core/v1"
)

const (
	// CheckNodeUnschedulablePred defines the name of predicate CheckNodeUnschedulablePredicate.
	CheckNodeUnschedulablePred = "CheckNodeUnschedulable"
)

// CheckNodeUnschedulablePredicate checks if a pod can be scheduled on a node with Unschedulable spec.
func CheckNodeUnschedulable(pod *v1.Pod, clusterInfo *resourceinfo.ClusterInfo) {

	for i, nodeInfo := range clusterInfo.Nodes {

		// If pod tolerate unschedulable taint, it's also tolerate `node.Spec.Unschedulable`.
		podToleratesUnschedulable := priorities.TolerationsTolerateTaint(pod.Spec.Tolerations, &v1.Taint{
			Key:    v1.TaintNodeUnschedulable,
			Effect: v1.TaintEffectNoSchedule,
		})

		// TODO (k82cn): deprecates `node.Spec.Unschedulable` in 1.13.
		if nodeInfo.Node.Spec.Unschedulable && !podToleratesUnschedulable {
			if i < len(clusterInfo.Nodes)-1 {
				clusterInfo.Nodes = append(clusterInfo.Nodes[:i], clusterInfo.Nodes[i+1:]...)
			} else {
				clusterInfo.Nodes = clusterInfo.Nodes[:i-1]
			}
		}
	}
}
