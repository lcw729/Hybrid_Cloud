package resourceinfo

import (
	cobrautil "Hybrid_Cloud/hybridctl/util"
	"context"
	"fmt"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

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

func NewNodeInfo(name string, pods ...*v1.Pod) *NodeInfo {
	ni := &NodeInfo{
		NodeName:           name,
		RequestedResources: &Resource{},
		AllocatableResources: &Resource{
			AllowedPodNumber: 1,
			MilliCPU:         5,
			Memory:           25,
		},
		ImageStates: make(map[string]*ImageStateSummary),
	}
	for _, pod := range pods {
		ni.AddPod(pod)
	}
	return ni
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

func (n *NodeInfo) VolumeLimits() map[v1.ResourceName]int64 {
	volumeLimits := map[v1.ResourceName]int64{}
	for k, v := range n.AllocatableResource().ScalarResources {
		if IsAttachableVolumeResourceName(k) {
			volumeLimits[k] = v
		}
	}
	return volumeLimits
}

func IsAttachableVolumeResourceName(name v1.ResourceName) bool {
	return strings.HasPrefix(string(name), v1.ResourceAttachableVolumesPrefix)
}

// AllocatableResource returns allocatable resources on a given node.
func (n *NodeInfo) AllocatableResource() Resource {
	if n == nil {
		return Resource{}
	}
	return *n.AllocatableResources
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
