package resourceinfo

import (
	cobrautil "Hybrid_Cloud/hybridctl/util"
	"Hybrid_Cloud/pkg/apis/hcpcluster/v1alpha1"
	hcpclusterv1alpha1 "Hybrid_Cloud/pkg/client/hcpcluster/v1alpha1/clientset/versioned"
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

type ClusterInfoList []*ClusterInfo

type ClusterInfo struct {
	ClusterName        string
	ClusterScore       int32
	Nodes              []*NodeInfo
	RequestedResources *Resources
	AllocableResources *Resources
	CapacityResources  *Resources
}

// CreateNodeInfoMap obtains a list of pods and pivots that list into a map where the keys are node names
// and the values are the aggregated information for that node.
func CreateClusterInfoMap(clusters *ClusterInfoList) map[string]*ClusterInfo {

	fmt.Println("[1-1] create clusterInfoMap")
	ClusterNameToInfo := make(map[string]*ClusterInfo)
	for _, cluster := range *clusters {
		clusterName := cluster.ClusterName
		if _, ok := ClusterNameToInfo[clusterName]; !ok {
			ClusterNameToInfo[clusterName] = cluster
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
		clusterInfo_list = append(clusterInfo_list, &clusterInfo)

	}

	return &clusterInfo_list
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
	nodeMetrics_list, err := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	fmt.Println(nodeMetrics_list.Items[0].Usage)

}

// ImageStates returns the state information of all images.
func (n *NodeInfo) GetImageStates() map[string]*ImageStateSummary {
	if n == nil {
		return nil
	}
	return n.ImageStates
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

// PodInfo is a wrapper to a Pod with additional pre-computed information to
// accelerate processing. This information is typically immutable (e.g., pre-processed
// inter-pod affinity selectors).
type PodInfo struct {
	ClusterName string
	NodeName    string
	PodName     string

	Pod *v1.Pod

	RequestedResources *Resources
	AllocableResources *Resources
	CapacityResources  *Resources

	RequiredAffinityTerms      []AffinityTerm
	RequiredAntiAffinityTerms  []AffinityTerm
	PreferredAffinityTerms     []WeightedAffinityTerm
	PreferredAntiAffinityTerms []WeightedAffinityTerm
	ParseError                 error
}

// AffinityTerm is a processed version of v1.PodAffinityTerm.
type AffinityTerm struct {
	Namespaces        sets.String
	Selector          labels.Selector
	TopologyKey       string
	NamespaceSelector labels.Selector
}

// WeightedAffinityTerm is a "processed" representation of v1.WeightedAffinityTerm.
type WeightedAffinityTerm struct {
	AffinityTerm
	Weight int32
}

// ImageStateSummary provides summarized information about the state of an image.
type ImageStateSummary struct {
	// Size of the image
	Size int64
	// Used to track how many nodes have this image
	NumNodes int
}

// NodeInfo is node level aggregated information.
type NodeInfo struct {
	ClusterName string
	NodeName    string
	// Overall node information.
	Node      *v1.Node
	NodeScore int32

	// Pods running on the node.
	Pods []*PodInfo

	// The subset of pods with affinity.
	PodsWithAffinity []*PodInfo

	// The subset of pods with required anti-affinity.
	PodsWithRequiredAntiAffinity []*PodInfo

	ImageStates map[string]*ImageStateSummary

	// Total requested resources of all pods on this node.
	RequestedResources   *Resource
	AllocatableResources *Resource
	CapacityResources    *Resource
}

// Resource is a collection of compute resource.
type Resource struct {
	MilliCPU         int64
	Memory           int64
	EphemeralStorage int64
	// We store allowedPodNumber (which is Node.Status.Allocatable.Pods().Value())
	// explicitly as int, to avoid conversions and improve performance.
	AllowedPodNumber int
	// ScalarResources
	ScalarResources map[v1.ResourceName]int64
}

type Resources struct {
	CPU     string
	Memory  string
	Fs      string
	Network string
}
