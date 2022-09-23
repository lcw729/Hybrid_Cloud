package test

import (
	// "github.com/KETI-Hybrid/hcp-scheduler-v1/pkg/algorithm/predicates"

	"testing"

	"github.com/KETI-Hybrid/hcp-scheduler-v1/src/resourceinfo"

	"k8s.io/klog"
)

func TestClusterUnschedulable(t *testing.T) {
	/*
		testdatas := []struct {
			node []*v1.Node
		}{
			{
				[]*v1.Node{
					NodeWithTaints("nodeA", []v1.Taint{{
						Key:    "foo",
						Value:  "bar",
						Effect: v1.TaintEffectPreferNoSchedule,
					}}),
					{
						Spec: v1.NodeSpec{
							Unschedulable: true,
						},
					},
				},
			},

			{ // the count of taints that are tolerated by pod, does not matter.

				[]*v1.Node{
					NodeWithTaints("nodeA", []v1.Taint{}),
					NodeWithTaints("nodeB", []v1.Taint{
						{
							Key:    "cpu-type",
							Value:  "arm64",
							Effect: v1.TaintEffectNoSchedule,
						},
					}),
					{
						Spec: v1.NodeSpec{
							Unschedulable: true,
						},
					},
				},
			},
			{ // the count of taints on a node that are not tolerated by pod, matters.

				[]*v1.Node{
					{
						Spec: v1.NodeSpec{
							Unschedulable: true,
						},
					},
					{
						Spec: v1.NodeSpec{
							Unschedulable: true,
						},
					},
				},
			},
			{ // taints-tolerations priority only takes care about the taints and tolerations that have effect PreferNoSchedule

				[]*v1.Node{
					NodeWithTaints("nodeA", []v1.Taint{}),
					NodeWithTaints("nodeB", []v1.Taint{
						{
							Key:    "cpu-type",
							Value:  "arm64",
							Effect: v1.TaintEffectNoSchedule,
						},
					}),
				},
			},
			{
				//PodWithTolerations([]v1.Toleration{}),
				[]*v1.Node{
					//Node without taints
					NodeWithTaints("nodeA", []v1.Taint{}),
					NodeWithTaints("nodeB", []v1.Taint{
						{
							Key:    "cpu-type",
							Value:  "arm64",
							Effect: v1.TaintEffectPreferNoSchedule,
						},
					}),
				},
			},
		}
		sched := scheduler.NewScheduler()
		sched.SchedulingResource = &v1.Pod{
			Spec: v1.PodSpec{
				Tolerations: []v1.Toleration{},
			},
		}
		for i, testdata := range testdatas {
			nodes_list := testdata.node
			cluster_name := "test_cluster" + strconv.Itoa(i+1)
			CreateTestClusters(&sched.ClusterInfoList, nodes_list, cluster_name)
		}

		sched.Filtering("CheckNodeUnschedulable")
	*/
}

func TestNodeUnschedulable(t *testing.T) {

	//pod := NewResourcePod(resourceinfo.Resource{MilliCPU: 1, Memory: 1})

	nodeList := []*resourceinfo.NodeInfo{
		resourceinfo.NewNodeInfo("node1", NewResourcePod(resourceinfo.Resource{MilliCPU: 10, Memory: 20})),
		resourceinfo.NewNodeInfo("node2", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 5})),
		resourceinfo.NewNodeInfo("node3", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 19})),
		resourceinfo.NewNodeInfo("node4", NewResourcePod(resourceinfo.Resource{MilliCPU: 5, Memory: 19})),
	}

	var clusterinfo resourceinfo.ClusterInfo
	clusterinfo.ClusterName = "test-cluster"
	clusterinfo.Nodes = append(clusterinfo.Nodes, nodeList...)

	klog.Infoln("===before NodeUnschedulable Filtering===")
	//predicates.CheckNodeUnschedulable(pod, &clusterinfo)
	klog.Infoln("===after NodeUnschedulable Filtering===")
	for _, node := range clusterinfo.Nodes {
		klog.Infoln((*node).NodeName)
	}

}
