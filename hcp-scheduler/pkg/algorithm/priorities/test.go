package priorities

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"strconv"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTestClusters(clusterinfo_list *resourceinfo.ClusterInfoList, node_list []*v1.Node, cluster_name string) {

	var cluster_info resourceinfo.ClusterInfo
	cluster_info.ClusterName = cluster_name

	for i := 0; i < len(node_list); i++ {
		var new_node resourceinfo.NodeInfo
		new_node.ClusterName = cluster_info.ClusterName
		new_node.NodeName = node_list[i].GetObjectMeta().GetName()
		new_node.Node = node_list[i]
		cluster_info.Nodes = append(cluster_info.Nodes, &new_node)
	}

	(*clusterinfo_list) = append((*clusterinfo_list), &cluster_info)
}

func CreateTestClusterNodeAffinity(clusterinfo_list *resourceinfo.ClusterInfoList) {
	// Test
	label1 := map[string]string{"foo": "bar"}                              // weight 2
	label2 := map[string]string{"key": "value"}                            // weight 4
	label3 := map[string]string{"az": "az1"}                               // weight 5
	label4 := map[string]string{"abc": "az11", "def": "az22"}              // weight 0
	label5 := map[string]string{"foo": "bar", "key": "value", "az": "az1"} // weight 2 + 4 + 5 = 11

	node_list_1 := []*v1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "machine1", Labels: label5}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine2", Labels: label2}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine3", Labels: label3}},
	}

	CreateTestClusters(clusterinfo_list, node_list_1, "test_cluster_1")

	node_list_2 := []*v1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "machine4", Labels: label1}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine5", Labels: label4}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine6", Labels: label3}},
	}
	CreateTestClusters(clusterinfo_list, node_list_2, "test_cluster_2")

	/*Test - NodeAffinity*/
	/*
		affinity1 := &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{{
					Weight: 2,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{{
							Key:      "foo",
							Operator: v1.NodeSelectorOpIn,
							Values:   []string{"bar"},
						}},
					},
				}},
			},
		}
		_ = affinity1
		affinity2 := &v1.Affinity{
			NodeAffinity: &v1.NodeAffinity{
				PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{
					{
						Weight: 2,
						Preference: v1.NodeSelectorTerm{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      "foo",
									Operator: v1.NodeSelectorOpIn,
									Values:   []string{"bar"},
								},
							},
						},
					},
					{
						Weight: 4,
						Preference: v1.NodeSelectorTerm{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      "key",
									Operator: v1.NodeSelectorOpIn,
									Values:   []string{"value"},
								},
							},
						},
					},
					{
						Weight: 5,
						Preference: v1.NodeSelectorTerm{
							MatchExpressions: []v1.NodeSelectorRequirement{
								{
									Key:      "az",
									Operator: v1.NodeSelectorOpIn,
									Values:   []string{"az1"},
								},
							},
						},
					},
				},
			},
		}

		priorities.CreateTestClusterNodeAffinity(&sched.ClusterInfoList)

		var rep int32 = 2

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
							Affinity: affinity2,
						},
					},
				},
			},
		}

		replicas := *test_deployment.Spec.RealDeploymentSpec.Replicas
	*/
}

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

// This function will create a set of nodes and pods and test the priority
// Nodes with zero,one,two,three,four and hundred taints are created
// Pods with zero,one,two,three,four and hundred tolerations are created

func CreateTestClusterTaintAndToleration(clusterinfo_list *resourceinfo.ClusterInfoList) {
	tests := []struct {
		pod  *v1.Pod
		node []*v1.Node
	}{
		{
			PodWithTolerations([]v1.Toleration{{
				Key:      "foo",
				Operator: v1.TolerationOpEqual,
				Value:    "bar",
				Effect:   v1.TaintEffectPreferNoSchedule,
			}}),
			[]*v1.Node{
				NodeWithTaints("nodeA", []v1.Taint{{
					Key:    "foo",
					Value:  "bar",
					Effect: v1.TaintEffectPreferNoSchedule,
				}}),
				NodeWithTaints("nodeB", []v1.Taint{
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
			PodWithTolerations([]v1.Toleration{
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
				NodeWithTaints("nodeA", []v1.Taint{}),
				NodeWithTaints("nodeB", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectPreferNoSchedule,
					},
				}),
				NodeWithTaints("nodeC", []v1.Taint{
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
			PodWithTolerations([]v1.Toleration{{
				Key:      "foo",
				Operator: v1.TolerationOpEqual,
				Value:    "bar",
				Effect:   v1.TaintEffectPreferNoSchedule,
			}}),
			[]*v1.Node{
				NodeWithTaints("nodeA", []v1.Taint{}),
				NodeWithTaints("nodeB", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectPreferNoSchedule,
					},
				}),
				NodeWithTaints("nodeC", []v1.Taint{
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
			PodWithTolerations([]v1.Toleration{
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
				NodeWithTaints("nodeA", []v1.Taint{}),
				NodeWithTaints("nodeB", []v1.Taint{
					{
						Key:    "cpu-type",
						Value:  "arm64",
						Effect: v1.TaintEffectNoSchedule,
					},
				}),
				NodeWithTaints("nodeC", []v1.Taint{
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
			PodWithTolerations([]v1.Toleration{}),
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

	for i, test := range tests {
		nodes_list := test.node
		cluster_name := "test_cluster" + strconv.Itoa(i+1)
		CreateTestClusters(clusterinfo_list, nodes_list, cluster_name)
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

/*
func TestNodeAffinityPriority(t *testing.T) {
	// Test
	label1 := map[string]string{"foo": "bar"}                              // weight 2
	label2 := map[string]string{"key": "value"}                            // weight 4
	label3 := map[string]string{"az": "az1"}                               // weight 5
	label4 := map[string]string{"abc": "az11", "def": "az22"}              // weight 0
	label5 := map[string]string{"foo": "bar", "key": "value", "az": "az1"} // weight 2 + 4 + 5 = 11
	fmt.Println("AAAA")
	var clusterinfo_list resourceinfo.ClusterInfoList
	node_list_1 := []*v1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "machine1", Labels: label5}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine2", Labels: label2}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine3", Labels: label3}},
	}
	fmt.Println("AAAA")
	CreateTestClusters(&clusterinfo_list, node_list_1, "test_cluster_1")

	node_list_2 := []*v1.Node{
		{ObjectMeta: metav1.ObjectMeta{Name: "machine4", Labels: label1}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine5", Labels: label4}},
		{ObjectMeta: metav1.ObjectMeta{Name: "machine6", Labels: label3}},
	}
	CreateTestClusters(&clusterinfo_list, node_list_2, "test_cluster_2")

	affinity1 := &v1.Affinity{
		NodeAffinity: &v1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{{
				Weight: 2,
				Preference: v1.NodeSelectorTerm{
					MatchExpressions: []v1.NodeSelectorRequirement{{
						Key:      "foo",
						Operator: v1.NodeSelectorOpIn,
						Values:   []string{"bar"},
					}},
				},
			}},
		},
	}
	_ = affinity1
	affinity2 := &v1.Affinity{
		NodeAffinity: &v1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{
				{
					Weight: 2,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "foo",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"bar"},
							},
						},
					},
				},
				{
					Weight: 4,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "key",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"value"},
							},
						},
					},
				},
				{
					Weight: 5,
					Preference: v1.NodeSelectorTerm{
						MatchExpressions: []v1.NodeSelectorRequirement{
							{
								Key:      "az",
								Operator: v1.NodeSelectorOpIn,
								Values:   []string{"az1"},
							},
						},
					},
				},
			},
		},
	}

	var rep int32 = 2

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
						Affinity: affinity2,
					},
				},
			},
		},
	}

	replicas := *test_deployment.Spec.RealDeploymentSpec.Replicas

	fake_pod := newPodFromHCPDeployment(test_deployment)

	for i := 0; i < int(replicas); i++ {
		scoring(fake_pod, &clusterinfo_list, "Affinity")
	}
}
*/
