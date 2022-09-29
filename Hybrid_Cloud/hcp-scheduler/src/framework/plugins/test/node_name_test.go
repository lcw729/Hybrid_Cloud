package test

import (
	"testing"

	"github.com/KETI-Hybrid/hcp-scheduler-v1/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func TestNodeName(t *testing.T) {
	test_pods := []*v1.Pod{
		{
			Spec: v1.PodSpec{
				NodeName: "foo",
			},
		},
		{
			Spec: v1.PodSpec{
				NodeName: "bar",
			},
		},
	}

	nodeList := []*resourceinfo.NodeInfo{
		{Node: &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "foo",
			},
		},
		},
	}

	var clusterinfo resourceinfo.ClusterInfo
	clusterinfo.ClusterName = "test-cluster"
	clusterinfo.Nodes = append(clusterinfo.Nodes, nodeList...)

	for _, pod := range test_pods {
		klog.Infoln(pod.Spec.NodeName)
		klog.Infoln("===before NodeName Filtering===")
		//framework.HCPFilterPlugin.Filter(pod, clusterinfo)
		klog.Infoln("===after NodeName Filtering===")
		for _, node := range clusterinfo.Nodes {
			klog.Infoln((*node).NodeName)
		}
	}
}
