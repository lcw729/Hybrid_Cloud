package priorities

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"

	v1 "k8s.io/api/core/v1"
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
