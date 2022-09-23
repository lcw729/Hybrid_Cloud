package test

import (
	"github.com/KETI-Hybrid/hcp-scheduler-v1/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// The two thresholds are used as bounds for the image score range. They correspond to a reasonable size range for
// container images compressed and stored in registries; 90%ile of images on dockerhub drops into this range.
const (
	mb           int64 = 1024 * 1024
	minThreshold int64 = 23 * mb
	maxThreshold int64 = 1000 * mb
)

func CreateTestClusterImageLocality(clusterinfo_list *resourceinfo.ClusterInfoList) {
	node403002000 := v1.NodeStatus{
		Images: []v1.ContainerImage{
			{
				Names: []string{
					"gcr.io/40:" + "latest",
					"gcr.io/40:v1",
					"gcr.io/40:v1",
				},
				SizeBytes: int64(40 * mb),
			},
			{
				Names: []string{
					"gcr.io/300:" + "latest",
					"gcr.io/300:v1",
				},
				SizeBytes: int64(300 * mb),
			},
			{
				Names: []string{
					"gcr.io/2000:" + "latest",
				},
				SizeBytes: int64(2000 * mb),
			},
		},
	}

	node25010 := v1.NodeStatus{
		Images: []v1.ContainerImage{
			{
				Names: []string{
					"gcr.io/250:" + "latest",
				},
				SizeBytes: int64(250 * mb),
			},
			{
				Names: []string{
					"gcr.io/10:" + "latest",
					"gcr.io/10:v1",
				},
				SizeBytes: int64(10 * mb),
			},
		},
	}

	nodeWithNoImages := v1.NodeStatus{}

	node_list_1 := []*v1.Node{makeImageNode("machine1", node403002000), makeImageNode("machine2", node25010)}
	CreateTestClusters(clusterinfo_list, node_list_1, "test_cluster_1")

	node_list_2 := []*v1.Node{makeImageNode("machine1", node403002000), makeImageNode("machine2", node25010), makeImageNode("machine3", nodeWithNoImages)}
	CreateTestClusters(clusterinfo_list, node_list_2, "test_cluster_2")
}

var TestPodsImageLocality []*v1.Pod = []*v1.Pod{
	{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image: "gcr.io/40",
				},
				{
					Image: "gcr.io/250",
				},
			},
		},
	},
	{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image: "gcr.io/40",
				},
				{
					Image: "gcr.io/300",
				},
			},
		},
	},
	{
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Image: "gcr.io/10",
				},
				{
					Image: "gcr.io/2000",
				},
			},
		},
	},
}

func makeImageNode(node string, status v1.NodeStatus) *v1.Node {
	return &v1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: node},
		Status:     status,
	}
}
