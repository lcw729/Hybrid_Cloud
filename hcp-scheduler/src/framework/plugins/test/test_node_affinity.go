package test

import (
	"github.com/KETI-Hybrid/hcp-scheduler-v1/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
