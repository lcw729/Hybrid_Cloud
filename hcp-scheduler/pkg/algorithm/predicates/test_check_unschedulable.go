package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/test"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"strconv"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NodeWithTaints(nodeName string, taints []v1.Taint) *v1.Node {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName,
		},
		Spec: v1.NodeSpec{
			Taints: taints,
		},
	}
}

func PodWithTolerations(tolerations []v1.Toleration) *v1.Pod {
	return &v1.Pod{
		Spec: v1.PodSpec{
			Tolerations: tolerations,
		},
	}
}

func CreateTestClusterUnschedulable(clusterinfo_list *resourceinfo.ClusterInfoList) {
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

	for i, testdata := range testdatas {
		nodes_list := testdata.node
		cluster_name := "test_cluster" + strconv.Itoa(i+1)
		test.CreateTestClusters(clusterinfo_list, nodes_list, cluster_name)
	}

	/*
		sched.SchedulingResource = &v1.Pod{
			Spec: v1.PodSpec{
				Tolerations: []v1.Toleration{},
			},
		}

		predicates.CreateTestClusterUnschedulable(&sched.ClusterInfoList)
		sched.Filtering("CheckNodeUnschedulable")
	*/
}
