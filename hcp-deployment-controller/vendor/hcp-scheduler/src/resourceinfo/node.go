package resourceinfo

import (
	"context"
	"strings"

	"hcp-pkg/util/clusterManager"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
	metrics "k8s.io/metrics/pkg/client/clientset/versioned"
)

var cm, _ = clusterManager.NewClusterManager()

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
// DefaultBindAllHostIP defines the default ip address used to bind to all host.
const DefaultBindAllHostIP = "0.0.0.0"

func NewNodeInfo(name string, pods ...*v1.Pod) *NodeInfo {
	r := new(Resource)
	ni := &NodeInfo{
		NodeName:           name,
		RequestedResources: *r,
		AllocatableResources: &Resource{
			AllowedPodNumber: 3,
			MilliCPU:         30,
			Memory:           25,
		},
		UsedPorts:   make(HostPortInfo),
		ImageStates: make(map[string]*ImageStateSummary),
	}
	(*ni).Node = &v1.Node{
		Spec: v1.NodeSpec{
			Taints: []v1.Taint{
				{
					Key:    v1.TaintNodeUnschedulable,
					Effect: v1.TaintEffectNoSchedule,
				},
			},
		},
	}

	for _, pod := range pods {
		ni.AddPod(pod)
	}
	return ni
}

func GetNodeInfo(clusterName string) []*NodeInfo {
	var nodeInfo []*NodeInfo
	config := cm.Cluster_configs[clusterName]
	cluster_client := kubernetes.NewForConfigOrDie(config)

	nodes, _ := cluster_client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	pods, _ := cluster_client.CoreV1().Pods(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})

	for _, node := range nodes.Items {
		ni := &NodeInfo{
			ClusterName: clusterName,
			NodeName:    node.Name,
			Node:        &node,
			UsedPorts:   make(HostPortInfo),
			ImageStates: make(map[string]*ImageStateSummary),
		}

		for _, pod := range pods.Items {
			ni.AddPod(&pod)
		}

		nodeInfo = append(nodeInfo, ni)
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
	config := cm.Cluster_configs[clusterName]

	mc, err := metrics.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// mc.MetricsV1beta1().NodeMetricses().Get(cotex"your node name", metav1.GetOptions{})
	nodeMetrics_list, _ := mc.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	klog.Infoln(nodeMetrics_list.Items[0].Usage)

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

// Add adds (ip, protocol, port) to HostPortInfo
func (h HostPortInfo) Add(ip, protocol string, port int32) {
	if port <= 0 {
		return
	}

	h.sanitize(&ip, &protocol)

	pp := NewProtocolPort(protocol, port)
	if _, ok := h[ip]; !ok {
		h[ip] = map[ProtocolPort]struct{}{
			*pp: {},
		}
		return
	}

	h[ip][*pp] = struct{}{}
}

// Remove removes (ip, protocol, port) from HostPortInfo
func (h HostPortInfo) Remove(ip, protocol string, port int32) {
	if port <= 0 {
		return
	}

	h.sanitize(&ip, &protocol)

	pp := NewProtocolPort(protocol, port)
	if m, ok := h[ip]; ok {
		delete(m, *pp)
		if len(h[ip]) == 0 {
			delete(h, ip)
		}
	}
}

// CheckConflict checks if the input (ip, protocol, port) conflicts with the existing
// ones in HostPortInfo.
func (h HostPortInfo) CheckConflict(ip, protocol string, port int32) bool {
	if port <= 0 {
		return false
	}

	h.sanitize(&ip, &protocol)

	pp := NewProtocolPort(protocol, port)

	// If ip is 0.0.0.0 check all IP's (protocol, port) pair
	if ip == DefaultBindAllHostIP {
		for _, m := range h {
			if _, ok := m[*pp]; ok {
				return true
			}
		}
		return false
	}

	// If ip isn't 0.0.0.0, only check IP and 0.0.0.0's (protocol, port) pair
	for _, key := range []string{DefaultBindAllHostIP, ip} {
		if m, ok := h[key]; ok {
			if _, ok2 := m[*pp]; ok2 {
				return true
			}
		}
	}

	return false
}

// NewProtocolPort creates a ProtocolPort instance.
func NewProtocolPort(protocol string, port int32) *ProtocolPort {
	pp := &ProtocolPort{
		Protocol: protocol,
		Port:     port,
	}

	if len(pp.Protocol) == 0 {
		pp.Protocol = string(v1.ProtocolTCP)
	}

	return pp
}

// sanitize the parameters
func (h HostPortInfo) sanitize(ip, protocol *string) {
	if len(*ip) == 0 {
		*ip = DefaultBindAllHostIP
	}
	if len(*protocol) == 0 {
		*protocol = string(v1.ProtocolTCP)
	}
}

// //
// func (n *NodeInfo) FilterNode() {
// 	n.IsFiltered = true
// }
