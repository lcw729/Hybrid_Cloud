package test

import (
	//"hcp-scheduler/pkg/algorithm/predicates"
	"testing"

	"hcp-scheduler/src/resourceinfo"

	"k8s.io/klog"
)

func TestNodeResourceFit(t *testing.T) {

	//	pod := NewResourcePod(resourceinfo.Resource{MilliCPU: 1, Memory: 1})

	nodeList := []*resourceinfo.NodeInfo{
		resourceinfo.NewNodeInfo("node1", NewResourcePod(resourceinfo.Resource{MilliCPU: 10, Memory: 20})),
		resourceinfo.NewNodeInfo("node2", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 5})),
		resourceinfo.NewNodeInfo("node3", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 19})),
		resourceinfo.NewNodeInfo("node4", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 19})),
	}

	var clusterinfo resourceinfo.ClusterInfo
	clusterinfo.ClusterName = "test-cluster"
	clusterinfo.Nodes = append(clusterinfo.Nodes, nodeList...)

	klog.Infoln("===before NodeResourceFit Filtering===")
	//predicates.NodeResourcesFit(pod, &clusterinfo)
	klog.Infoln("===after NodeResourceFit Filtering===")
	for _, node := range clusterinfo.Nodes {
		klog.Infoln((*node).NodeName)
	}

	klog.Infoln("===before NodeUnschedulable Filtering===")
	//	predicates.CheckNodeUnschedulable(pod, &clusterinfo)
	klog.Infoln("===after NodeUnschedulable Filtering===")
	for _, node := range clusterinfo.Nodes {
		klog.Infoln((*node).NodeName)
	}

}
