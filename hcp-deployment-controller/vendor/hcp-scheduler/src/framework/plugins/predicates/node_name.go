package predicates

import (
	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type NodeName struct{}

func (pl *NodeName) Name() string {
	return plugins.NodeName
}

// Filter invoked at the filter extension point.
func (pl *NodeName) Filter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	for _, nodeInfo := range clusterInfo.Nodes {
		if nodeInfo.Node == nil {
			klog.Infoln("node not found")
			//nodeInfo.FilterNode()
			//clusterInfo.MinusOneAvailableNodes()
			continue
		}

		klog.Info("fits ", Fits(pod, nodeInfo))
		klog.Info("podname ", pod.Spec.NodeName)
		klog.Info("nodename ", nodeInfo.NodeName)
		if Fits(pod, nodeInfo) {
			return false
			//nodeInfo.FilterNode()
			//clusterInfo.MinusOneAvailableNodes()
		} else {
			klog.Infoln("Node Name is unmatched")
		}
	}
	return true
}

// Fits actually checks if the pod fits the node.
func Fits(pod *v1.Pod, nodeInfo *resourceinfo.NodeInfo) bool {
	return len(pod.Spec.NodeName) == 0 || pod.Spec.NodeName == nodeInfo.NodeName
}
