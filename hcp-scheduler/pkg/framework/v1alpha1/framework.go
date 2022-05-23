package v1alpha1

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

const (
	NodeName          = "NodeName"
	NodePorts         = "NodePorts"
	NodeResourcesFit  = "NodeResourcesFit"
	NodeUnschedulable = "NodeUnschedulable"
	BalanceAllocation = "BalancedAllocation"
	ImageLocality     = "ImageLocality"
	TaintToleration   = "TaintToleration"
)

type hcpFramework struct {
	filterPlugins []HCPFilterPlugin
	scorePlugins  []HCPScorePlugin
}

func NewFramework() hcpFramework {
	framework := &hcpFramework{
		filterPlugins: []HCPFilterPlugin{},
		scorePlugins:  []HCPScorePlugin{},
	}
	return *framework
}

func (f *hcpFramework) RunFilterPluginsOnClusters(pod *v1.Pod, status *CycleStatus, clusterInfoList *resourceinfo.ClusterInfoList) {
	result := make(map[string]bool)
	var isFiltered bool
	for _, cluster := range *clusterInfoList {
		isFiltered = true
		for _, plugin := range f.filterPlugins {
			fmt.Println(plugin.Name())
			isFiltered = plugin.Filter(pod, status, &cluster)
			/*
			  result : true => pass
			  result : false => fail
			*/
			result[cluster.ClusterName] = !isFiltered
			// 하나의 plugin이라도 fail이면 다음 클러스터 필터링 시작
			if !result[cluster.ClusterName] {
				break
			}
		}
	}
}

func (f *hcpFramework) RunScorePluginsOnClusters(pod *v1.Pod, status *CycleStatus, clusterInfoList *resourceinfo.ClusterInfoList) {
	result := make(map[string]int64)
	var score int64
	for _, cluster := range *clusterInfoList {
		score = 0
		for _, plugin := range f.scorePlugins {
			fmt.Println(plugin.Name())
			score = plugin.Score(pod, status, &cluster)
			result[cluster.ClusterName] += score
		}
	}
}
