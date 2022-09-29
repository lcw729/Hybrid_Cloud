package resourceinfo

import (
	"context"

	"hcp-pkg/apis/hcpcluster/v1alpha1"
	hcpclusterv1alpha1 "hcp-pkg/client/hcpcluster/v1alpha1/clientset/versioned"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
)

func NewClusterInfoList() *ClusterInfoList {

	var clusterInfo_list ClusterInfoList
	joinCluster_list, err := JoinClusterList()
	if err != nil {
		klog.Errorln(err)
		return nil
	}

	for _, hcpcluster := range joinCluster_list {
		cluster_name := hcpcluster.Name
		clusterInfo := &ClusterInfo{
			ClusterName:  cluster_name,
			ClusterScore: 0,
			Nodes:        GetNodeInfo(cluster_name),
		}
		clusterInfo_list = append(clusterInfo_list, clusterInfo)
	}

	return &clusterInfo_list
}

// CreateNodeInfoMap obtains a list of pods and pivots that list into a map where the keys are node names
// and the values are the aggregated information for that node.
func CreateClusterInfoMap(clusters *ClusterInfoList) map[string]*ClusterInfo {
	ClusterNameToInfo := make(map[string]*ClusterInfo, len(*clusters))
	for _, cluster := range *clusters {
		clusterName := cluster.ClusterName
		if _, ok := ClusterNameToInfo[clusterName]; !ok {
			ClusterNameToInfo[clusterName] = cluster
		}
	}

	return ClusterNameToInfo
}

func JoinClusterList() ([]v1alpha1.HCPCluster, error) {
	var joinCluster_list []v1alpha1.HCPCluster
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Errorln(err)
		return nil, err
	}

	cluster_client := hcpclusterv1alpha1.NewForConfigOrDie(config)

	cluster_list, err := cluster_client.HcpV1alpha1().HCPClusters("hcp").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln(err)
		return nil, err
	}

	for _, hcpcluster := range cluster_list.Items {
		if hcpcluster.Spec.JoinStatus == "JOIN" {
			joinCluster_list = append(joinCluster_list, hcpcluster)
		}
	}

	return joinCluster_list, nil
}

// createImageExistenceMap returns a map recording on which nodes the images exist, keyed by the images' names.
func CreateImageExistenceMap(clusterinfoList *ClusterInfoList) map[string]sets.String {
	imageExistenceMap := make(map[string]sets.String)

	for _, cluster := range *clusterinfoList {
		nodes := cluster.Nodes
		for _, node := range nodes {
			for _, image := range node.Node.Status.Images {
				for _, name := range image.Names {
					if _, ok := imageExistenceMap[name]; !ok {
						imageExistenceMap[name] = sets.NewString(node.Node.Name)
					} else {
						imageExistenceMap[name].Insert(node.Node.Name)
					}
				}
			}
		}
	}
	klog.Infoln()
	return imageExistenceMap
}
