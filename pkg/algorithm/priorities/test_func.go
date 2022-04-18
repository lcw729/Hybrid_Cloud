package priorities

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/internal/scoretable"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var algorithmMap = map[string]func(*v1.Pod, *v1.Node) int32{

	"Affinity":        NodeAffinity,
	"TaintToleration": TaintToleration,
}

func scoring(pod *v1.Pod, clusterinfoList *resourceinfo.ClusterInfoList, algorithm string) {

	var score int32
	score_table := scoretable.NewScoreTable(len(*clusterinfoList))

	for _, clusterinfo := range *clusterinfoList {
		fmt.Println("==>", clusterinfo.ClusterName)
		score = 0
		node_list := (*clusterinfo).Nodes

		for _, node := range node_list {
			node_score := algorithmMap[algorithm](pod, node.Node)
			if node_score == -1 {
				fmt.Println("fail to scoring node")
				return
			} else {
				fmt.Println(node.NodeName, "score :", node_score)
				score += node_score
			}
		}
		fmt.Println("*", clusterinfo.ClusterName, "total score :", score)
		(*score_table)[clusterinfo.ClusterName] = float32(score)
	}
}

func createTestClusters(clusterinfo_list *resourceinfo.ClusterInfoList, node_list []*v1.Node, cluster_name string) {

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

func printSchedulingResult(targets []v1alpha1.Target) {
	fmt.Println("========scheduling result========")
	for _, i := range targets {
		fmt.Println("target cluster :", i.Cluster)
		fmt.Println("replicas       :", *i.Replicas)
		fmt.Println()
	}
}

func newPodFromHCPDeployment(deployment *v1alpha1.HCPDeployment) *v1.Pod {

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deployment.GetObjectMeta().GetName() + "-pod",
			Annotations: deployment.Annotations,
			Labels:      deployment.Labels,
		},
		Spec: deployment.Spec.RealDeploymentSpec.Template.Spec,
	}

	return pod
}

func updateSchedulingResult(scheduling_result *[]v1alpha1.Target, cluster string) bool {

	for i, target := range *scheduling_result {
		// 이미 target cluster 목록에 cluster가 있는 경우
		if target.Cluster == cluster {
			// replicas 개수 증가
			temp := *target.Replicas
			temp += 1
			target.Replicas = &temp
			(*scheduling_result)[i] = target
			return true
		}
	}

	// target cluster 목록에 cluster가 없는 경우

	// replicas 개수 1로 설정
	var new_target v1alpha1.Target
	new_target.Cluster = cluster
	var one int32 = 1
	new_target.Replicas = &one
	*scheduling_result = append((*scheduling_result), new_target)
	return true
}
