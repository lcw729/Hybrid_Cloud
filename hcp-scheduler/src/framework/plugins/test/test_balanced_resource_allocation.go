package test

import (
	"github.com/KETI-Hybrid/hcp-scheduler-v1/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateTestClusterNodeResourcesBalancedAllocation(clusterinfo_list *resourceinfo.ClusterInfoList) {
	node_list_1 := []*v1.Node{makeNode("machine1", 4000, 10000), makeNode("machine2", 4000, 10000)}
	CreateTestClusters(clusterinfo_list, node_list_1, "test_cluster_1")

	node_list_2 := []*v1.Node{makeNode("machine1", 4000, 10000), makeNode("machine2", 6000, 10000)}
	CreateTestClusters(clusterinfo_list, node_list_2, "test_cluster_2")

	node_list_3 := []*v1.Node{makeNode("machine1", 0, 0), makeNode("machine2", 0, 0)}
	CreateTestClusters(clusterinfo_list, node_list_3, "test_cluster_3")
}

var TestPodsNodeResourcesBalancedAllocation []*v1.Pod = []*v1.Pod{
	{Spec: v1.PodSpec{
		Containers: []v1.Container{
			{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("1000m"),
						v1.ResourceMemory: resource.MustParse("2000"),
					},
				},
			},
			{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("2000m"),
						v1.ResourceMemory: resource.MustParse("3000"),
					},
				},
			},
		}},
	},
	{Spec: v1.PodSpec{
		Containers: []v1.Container{
			{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("0m"),
						v1.ResourceMemory: resource.MustParse("0"),
					},
				},
			},
			{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("0m"),
						v1.ResourceMemory: resource.MustParse("0"),
					},
				},
			},
		},
		Volumes: []v1.Volume{
			{
				VolumeSource: v1.VolumeSource{
					AWSElasticBlockStore: &v1.AWSElasticBlockStoreVolumeSource{VolumeID: "ovp1"},
				},
			},
		}},
	},
	{Spec: v1.PodSpec{
		NodeName: "machine3",
		Containers: []v1.Container{
			{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("1000m"),
						v1.ResourceMemory: resource.MustParse("2000"),
					},
				},
			},
			{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("2000m"),
						v1.ResourceMemory: resource.MustParse("3000"),
					},
				},
			},
		},
	}},
	{Spec: v1.PodSpec{
		NodeName: "machine1",
		Containers: []v1.Container{
			{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("1000m"),
						v1.ResourceMemory: resource.MustParse("0"),
					},
				},
			},
			{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceCPU:    resource.MustParse("2000m"),
						v1.ResourceMemory: resource.MustParse("0"),
					},
				},
			},
		},
	}},
}

func makeNode(node string, milliCPU, memory int64) *v1.Node {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: node},
		Status: v1.NodeStatus{
			Capacity: v1.ResourceList{
				v1.ResourceCPU:    *resource.NewMilliQuantity(milliCPU, resource.DecimalSI),
				v1.ResourceMemory: *resource.NewQuantity(memory, resource.BinarySI),
				"pods":            *resource.NewQuantity(100, resource.DecimalSI),
			},
			Allocatable: v1.ResourceList{

				v1.ResourceCPU:    *resource.NewMilliQuantity(milliCPU, resource.DecimalSI),
				v1.ResourceMemory: *resource.NewQuantity(memory, resource.BinarySI),
				"pods":            *resource.NewQuantity(100, resource.DecimalSI),
			},
		},
	}
}
