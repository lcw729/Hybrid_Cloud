package test

import (
	"Hybrid_Cloud/hcp-scheduler/src/resourceinfo"
	"fmt"
	"testing"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestNodeName(t *testing.T) {
	test_pods := []*v1.Pod{
		{
			Spec: v1.PodSpec{
				NodeName: "foo",
			},
		},
		{
			Spec: v1.PodSpec{
				NodeName: "bar",
			},
		},
	}

	nodeList := []*resourceinfo.NodeInfo{
		{Node: &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
		},
		},
	}

	var clusterinfo resourceinfo.ClusterInfo
	clusterinfo.ClusterName = "test-cluster"
	clusterinfo.Nodes = append(clusterinfo.Nodes, nodeList...)

	for _, pod := range test_pods {
		fmt.Println(pod.Spec.NodeName)
		fmt.Println("===before NodeName Filtering===")
		//framework.HCPFilterPlugin.Filter(pod, clusterinfo)
		fmt.Println("===after NodeName Filtering===")
		for _, node := range clusterinfo.Nodes {
			fmt.Println((*node).NodeName)
		}
	}
}
