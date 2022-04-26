package test

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/predicates"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
)

func newResourcePod(usage ...resourceinfo.Resource) *v1.Pod {
	var containers []v1.Container
	for _, req := range usage {
		containers = append(containers, v1.Container{
			Resources: v1.ResourceRequirements{Requests: req.ResourceList()},
		})
	}
	return &v1.Pod{
		Spec: v1.PodSpec{
			Containers: containers,
		},
	}
}

func TestNodeResourceFit(t *testing.T) {

	pod := newResourcePod(resourceinfo.Resource{MilliCPU: 1, Memory: 1})

	nodeList := []*resourceinfo.NodeInfo{
		resourceinfo.NewNodeInfo("node1", newResourcePod(resourceinfo.Resource{MilliCPU: 10, Memory: 20})),
		resourceinfo.NewNodeInfo("node2", newResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 5})),
		resourceinfo.NewNodeInfo("node3", newResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 19})),
		resourceinfo.NewNodeInfo("node4", newResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 19})),
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

}
