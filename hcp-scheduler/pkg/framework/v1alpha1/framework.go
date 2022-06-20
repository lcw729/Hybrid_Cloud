package v1alpha1

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/framework/plugins/predicates"
	"Hybrid_Cloud/hcp-scheduler/pkg/framework/plugins/priorities"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

type hcpFramework struct {
	filterPlugins []HCPFilterPlugin
	scorePlugins  []HCPScorePlugin
}

func NewFramework() *hcpFramework {
	framework := &hcpFramework{
		filterPlugins: []HCPFilterPlugin{
			&predicates.NodeName{},
			&predicates.NodePorts{},
			&predicates.NodeResourcesFit{},
			&predicates.NodeUnschedulable{},
			&predicates.TaintToleration{},
			&predicates.JoinCheck{},
		},
		scorePlugins: []HCPScorePlugin{
			&priorities.BalanceAllocation{},
			&priorities.ImageLocality{},
			&priorities.NodeAffinity{},
		},
	}
	return framework
}

func (f *hcpFramework) RunFilterPluginsOnClusters(algorithms []string, pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfoList *resourceinfo.ClusterInfoList) {
	result := make(map[string]bool)
	var isFiltered bool

	var plugins []HCPFilterPlugin
	for _, i := range algorithms {
		plugin := f.stringTOHCPFilterPlugin(i)
		if plugin != nil {
			fmt.Println(plugin.Name())
			plugins = append(plugins, plugin)
		} else {
			fmt.Println(i, " : no such filter algorithm")
		}
	}

	for _, cluster := range *clusterInfoList {
		fmt.Println(">>", cluster.ClusterName)
		isFiltered = false
		for _, plugin := range plugins {
			fmt.Println("[plugin]", plugin.Name())
			isFiltered = plugin.Filter(pod, status, cluster)
			/*
			  result : true => 필터 O
			  result : false => 필터 X
			*/
			fmt.Println(">>>>>", isFiltered)
			(*cluster).IsFiltered = isFiltered
			result[cluster.ClusterName] = isFiltered
			//하나의 plugin이라도 true이면 다음 클러스터 필터링 시작
			if result[cluster.ClusterName] {
				break
			}
		}
	}
}

func (f *hcpFramework) stringTOHCPFilterPlugin(name string) HCPFilterPlugin {
	for _, p := range f.filterPlugins {
		if p.Name() == name {
			return p
		}
	}
	return nil
}

func (f *hcpFramework) RunScorePluginsOnClusters(algorithms []string, pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfoList *resourceinfo.ClusterInfoList) {
	result := make(map[string]int64)
	var score int64
	var plugins []HCPScorePlugin

	for _, i := range algorithms {
		plugin := f.stringTOHCPScorePlugin(i)
		if plugin != nil {
			fmt.Println(plugin.Name())
			plugins = append(plugins, plugin)
		} else {
			fmt.Println(i, " : no such filter algorithm")
		}
	}

	for _, cluster := range *clusterInfoList {
		score = 0
		fmt.Println(">>", cluster.ClusterName)
		if cluster.IsFiltered {
			fmt.Println(cluster.ClusterName, "is already filtered")
		} else {
			for _, plugin := range plugins {
				fmt.Println("[plugin]", plugin.Name())
				score = plugin.Score(pod, status, cluster)
				fmt.Println(score)
				result[cluster.ClusterName] += score
			}
		}
		(*cluster).ClusterScore = int32(result[cluster.ClusterName])
	}
}

func (f *hcpFramework) stringTOHCPScorePlugin(name string) HCPScorePlugin {
	for _, p := range f.scorePlugins {
		if p.Name() == name {
			return p
		}
	}
	return nil
}
