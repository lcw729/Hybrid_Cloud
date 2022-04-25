package priorities

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/test"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"strconv"

	v1 "k8s.io/api/core/v1"
)

// This function will create a set of nodes and pods and test the priority
// Nodes with zero,one,two,three,four and hundred taints are created
// Pods with zero,one,two,three,four and hundred tolerations are created

func CreateTestClusterTaintAndToleration(clusterinfo_list *resourceinfo.ClusterInfoList) {
	testdatas := []struct {
		pod  *v1.Pod
		node []*v1.Node
	}{
		{
			test.PodWithTolerations([]v1.Toleration{{
				Key:      "foo",
				Operator: v1.TolerationOpEqual,
				Value:    "bar",
				Effect:   v1.TaintEffectPreferNoSchedule,
			}}),
			[]*v1.Node{
				test.NodeWithTaints("nodeA", []v1.Taint{{
					Key:    "foo",
					Value:  "bar",
					Effect: v1.TaintEffectPreferNoSchedule,
				}}),
				test.NodeWithTaints("nodeB", []v1.Taint{
					{
						Key:    "foo",
						Value:  "bar",
						Effect: v1.TaintEffectPreferNoSchedule,
					}, {
						Key:    "foo",
						Value:  "blah",
						Effect: v1.TaintEffectPreferNoSchedule,
					}}),
			},
		},

		{ // the count of taints that are tolerated by pod, does not matter.
			test.PodWithTolerations([]v1.Toleration{
				// {
				// 	Key:      "cpu-type",
				// 	Operator: v1.TolerationOpEqual,
				// 	Value:    "arm64",
				// 	Effect:   v1.TaintEffectPreferNoSchedule,
				// },
				{
					Key:      "disk-type",
					Operator: v1.TolerationOpEqual,
					Value:    "ssd",
					Effect:   v1.TaintEffectPreferNoSchedule,
				},
			}),
			[]*v1.Node{
				test.NodeWithTaints("nodeA", []v1.Taint{}),
				test.NodeWithTaints("nodeB", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectPreferNoSchedule,
					},
				}),
				test.NodeWithTaints("nodeC", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectPreferNoSchedule,
					}, {
						Key:    "disk-type",
						Value:  "ssd",
						Effect: v1.TaintEffectPreferNoSchedule,
					},
				}),
			},
		},
		{ // the count of taints on a node that are not tolerated by pod, matters.
			test.PodWithTolerations([]v1.Toleration{{
				Key:      "foo",
				Operator: v1.TolerationOpEqual,
				Value:    "bar",
				Effect:   v1.TaintEffectPreferNoSchedule,
			}}),
			[]*v1.Node{
				test.NodeWithTaints("nodeA", []v1.Taint{}),
				test.NodeWithTaints("nodeB", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectPreferNoSchedule,
					},
				}),
				test.NodeWithTaints("nodeC", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectPreferNoSchedule,
					}, {
						Key:    "disk-type",
						Value:  "ssd",
						Effect: v1.TaintEffectPreferNoSchedule,
					},
				}),
			},
		},
		{ // taints-tolerations priority only takes care about the taints and tolerations that have effect PreferNoSchedule
			test.PodWithTolerations([]v1.Toleration{
				{
					Key:      "cpu-type",
					Operator: v1.TolerationOpEqual,
					Value:    "arm64",
					Effect:   v1.TaintEffectNoSchedule,
				}, {
					Key:      "disk-type",
					Operator: v1.TolerationOpEqual,
					Value:    "ssd",
					Effect:   v1.TaintEffectNoSchedule,
				},
			}),
			[]*v1.Node{
				test.NodeWithTaints("nodeA", []v1.Taint{}),
				test.NodeWithTaints("nodeB", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectNoSchedule,
					},
				}),
				test.NodeWithTaints("nodeC", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectPreferNoSchedule,
					}, {
						Key:    "disk-type",
						Value:  "ssd",
						Effect: v1.TaintEffectPreferNoSchedule,
					},
				}),
			},
		},
		{
			test.PodWithTolerations([]v1.Toleration{}),
			[]*v1.Node{
				//Node without taints
				test.NodeWithTaints("nodeA", []v1.Taint{}),
				test.NodeWithTaints("nodeB", []v1.Taint{
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
		var rep int32 = 2

		test := tests[2].pod
		fmt.Println("TaintToleration")
		test_deployment := &v1alpha1.HCPDeployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:        "test_deployment",
				Annotations: map[string]string{},
			},
			Spec: v1alpha1.HCPDeploymentSpec{
				RealDeploymentSpec: appsv1.DeploymentSpec{
					Replicas: &rep,
					Template: v1.PodTemplateSpec{
						Spec: v1.PodSpec{
							Tolerations: test.Spec.Tolerations,
						},
					},
				},
			},
		}

		replicas := *test_deployment.Spec.RealDeploymentSpec.Replicas

		fake_pod := newPodFromHCPDeployment(test_deployment)

		for i := 0; i < int(replicas); i++ {
			scoring(fake_pod, &clusterinfo_list, "TaintToleration")
		}
	*/
}
