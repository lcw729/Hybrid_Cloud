package resourceinfo

import (
	cobrautil "Hybrid_Cloud/hybridctl/util"
	"Hybrid_Cloud/pkg/apis/hcpcluster/v1alpha1"
	hcpclusterv1alpha1 "Hybrid_Cloud/pkg/client/hcpcluster/v1alpha1/clientset/versioned"
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

func NewClusterInfoList() *ClusterInfoList {

	var clusterInfo_list ClusterInfoList
	joinCluster_list, err := JoinClusterList()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for _, hcpcluster := range joinCluster_list {
		cluster_name := hcpcluster.Name
		clusterInfo := ClusterInfo{
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

	fmt.Println("[1-1] create clusterInfoMap")
	ClusterNameToInfo := make(map[string]*ClusterInfo)
	for _, cluster := range *clusters {
		clusterName := cluster.ClusterName
		if _, ok := ClusterNameToInfo[clusterName]; !ok {
			ClusterNameToInfo[clusterName] = &cluster
		}
	}

	return ClusterNameToInfo
}

/*
// createNodeInfoMap obtains a list of pods and pivots that list into a map
// where the keys are node names and the values are the aggregated information
// for that node.
func CreateNodeInfoMap(clusters *ClusterInfoList) map[string]*NodeInfo {
	nodeNameToInfo := make(map[string]*NodeInfo)
	//
	// 	for _, pod := range pods {
	// 		nodeName := pod.Spec.NodeName
	// 		if _, ok := nodeNameToInfo[nodeName]; !ok {
	// 			nodeNameToInfo[nodeName] = NewNodeInfo()
	// 		}
	// 		nodeNameToInfo[nodeName].AddPod(pod)
	// 	}
	//
	imageExistenceMap := CreateImageExistenceMap(clusters)

	for _, cluster := range *clusters {
		nodes := cluster.Nodes
		for _, node := range nodes {
			if _, ok := nodeNameToInfo[node.Node.Name]; !ok {
				nodeNameToInfo[node.Node.Name] = node
			}
			nodeInfo := nodeNameToInfo[node.Node.Name]
			nodeInfo.ImageStates = GetNodeImageStates(node.Node, imageExistenceMap)
			node = nodeInfo
		}
	}
	return nodeNameToInfo
}
*/
func JoinClusterList() ([]v1alpha1.HCPCluster, error) {

	var joinCluster_list []v1alpha1.HCPCluster
	config, err := cobrautil.BuildConfigFromFlags("kube-master", "/root/.kube/config")
	if err != nil {
		fmt.Println("err")
		return nil, err
	}

	cluster_client := hcpclusterv1alpha1.NewForConfigOrDie(config)

	cluster_list, err := cluster_client.HcpV1alpha1().HCPClusters("hcp").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println("err")
		return nil, err
	}

	for _, hcpcluster := range cluster_list.Items {
		if hcpcluster.Spec.JoinStatus == "JOIN" {
			joinCluster_list = append(joinCluster_list, hcpcluster)
		}
	}

	return joinCluster_list, nil
}

func GetNodeInfo(clusterName string) []*NodeInfo {
	var nodeInfo []*NodeInfo
	config, err := cobrautil.BuildConfigFromFlags(clusterName, "/root/.kube/config")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	cluster_client := kubernetes.NewForConfigOrDie(config)

	_, err = cluster_client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nodeInfo
}

func NodeMetrics(clusterName string) {
	config, err := cobrautil.BuildConfigFromFlags(clusterName, "/root/.kube/config")
	if err != nil {
		fmt.Println("this error")
	}

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// mc.MetricsV1beta1().NodeMetricses().Get(cotex"your node name", metav1.GetOptions{})
	nodeMetrics_list, _ := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	fmt.Println(nodeMetrics_list.Items[0].Usage)

}

// ImageStates returns the state information of all images.
func (n *NodeInfo) GetImageStates() map[string]*ImageStateSummary {
	if n == nil {
		return nil
	}
	return n.ImageStates
}

// getNodeImageStates returns the given node's image states based on the given imageExistence map.
func GetNodeImageStates(node *v1.Node, imageExistenceMap map[string]sets.String) map[string]*ImageStateSummary {
	imageStates := make(map[string]*ImageStateSummary)

	for _, image := range (*node).Status.Images {
		for _, name := range image.Names {
			imageStates[name] = &ImageStateSummary{
				Size:     image.SizeBytes,
				NumNodes: len(imageExistenceMap[name]),
			}
		}
	}
	return imageStates
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
	fmt.Println()
	return imageExistenceMap
}
