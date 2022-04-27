package test

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/predicates"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"fmt"
	"testing"
)

func TestNodeResourceFit(t *testing.T) {

	pod := NewResourcePod(resourceinfo.Resource{MilliCPU: 1, Memory: 1})

	nodeList := []*resourceinfo.NodeInfo{
		resourceinfo.NewNodeInfo("node1", NewResourcePod(resourceinfo.Resource{MilliCPU: 10, Memory: 20})),
		resourceinfo.NewNodeInfo("node2", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 5})),
		resourceinfo.NewNodeInfo("node3", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 19})),
		resourceinfo.NewNodeInfo("node4", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 19})),
	}

	var clusterinfo resourceinfo.ClusterInfo
	clusterinfo.ClusterName = "test-cluster"
	clusterinfo.Nodes = append(clusterinfo.Nodes, nodeList...)

	fmt.Println("===before NodeResourceFit Filtering===")
	predicates.NodeResourcesFit(pod, &clusterinfo)
	fmt.Println("===after NodeResourceFit Filtering===")
	for _, node := range clusterinfo.Nodes {
		fmt.Println((*node).NodeName)
	}

	fmt.Println("===before NodeUnschedulable Filtering===")
	predicates.CheckNodeUnschedulable(pod, &clusterinfo)
	fmt.Println("===after NodeUnschedulable Filtering===")
	for _, node := range clusterinfo.Nodes {
		fmt.Println((*node).NodeName)
	}

}
