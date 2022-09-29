package predicates

import (
	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type NodeUnschedulable struct{}

func (pl *NodeUnschedulable) Name() string {
	return plugins.NodeUnschedulable
}

// CheckNodeUnschedulablePredicate checks if a pod can be scheduled on a node with Unschedulable spec.
func (pl *NodeUnschedulable) Filter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	for _, nodeInfo := range clusterInfo.Nodes {
		// If pod tolerate unschedulable taint, it's also tolerate `node.Spec.Unschedulable`.
		podToleratesUnschedulable := TolerationsTolerateTaint(pod.Spec.Tolerations, &v1.Taint{
			Key:    v1.TaintNodeUnschedulable,
			Effect: v1.TaintEffectNoSchedule,
		})

		// TODO (k82cn): deprecates `node.Spec.Unschedulable` in 1.13.
		if nodeInfo.Node.Spec.Unschedulable && podToleratesUnschedulable {
			return false
		} else {
			klog.Infoln("Node has Unschedulable Taint.\nbut, Pod hasn't Unschedulable Toleration.")
		}
	}
	return true
}
