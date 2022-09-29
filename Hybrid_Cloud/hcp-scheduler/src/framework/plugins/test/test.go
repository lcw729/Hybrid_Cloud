package test

import (
	"github.com/KETI-Hybrid/hcp-scheduler-v1/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewResourcePod(usage ...resourceinfo.Resource) *v1.Pod {
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

func CreateTestClusters(clusterinfo_list *resourceinfo.ClusterInfoList, node_list []*v1.Node, cluster_name string) {

	var cluster_info resourceinfo.ClusterInfo
	cluster_info.ClusterName = cluster_name
	imageExistenceMap := resourceinfo.CreateImageExistenceMap(clusterinfo_list)
	// 테스트를 위해 clusterinfolist에 새로운 노드 등록하기
	for i := 0; i < len(node_list); i++ {
		var new_node resourceinfo.NodeInfo
		new_node.ClusterName = cluster_info.ClusterName
		new_node.NodeName = node_list[i].GetObjectMeta().GetName()
		new_node.Node = node_list[i]
		new_node.ImageStates = resourceinfo.GetNodeImageStates(node_list[i], imageExistenceMap)
		cluster_info.Nodes = append(cluster_info.Nodes, &new_node)
	}

	(*clusterinfo_list) = append((*clusterinfo_list), &cluster_info)
}
