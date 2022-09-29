package v1alpha1

import (
	"hcp-scheduler/src/framework/plugins/predicates"
	"hcp-scheduler/src/framework/plugins/priorities"
	"hcp-scheduler/src/resourceinfo"
	"hcp-scheduler/src/util"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type hcpFramework struct {
	pluginScoreList PluginScoreList
	filterPlugins   []HCPFilterPlugin
	scorePlugins    []HCPScorePlugin
}

func NewFramework() *hcpFramework {
	pluginScoreList := make(PluginScoreList)
	framework := &hcpFramework{
		pluginScoreList: pluginScoreList,
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
			klog.Infoln(plugin.Name())
			plugins = append(plugins, plugin)
		} else {
			klog.Infoln(i, " : no such filter algorithm")
		}
	}

	for _, cluster := range *clusterInfoList {
		klog.Infoln(">>", cluster.ClusterName)
		isFiltered = false
		for _, plugin := range plugins {
			klog.Infoln("[plugin]", plugin.Name())
			isFiltered = plugin.Filter(pod, status, cluster)
			/*
			  result : true => 필터 O
			  result : false => 필터 X
			*/
			klog.Infoln(">>>>>", isFiltered)
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
	// result := new(map[string]int64)
	var score int64
	var plugins []HCPScorePlugin

	for _, i := range algorithms {
		plugin := f.stringTOHCPScorePlugin(i)
		if plugin != nil {
			klog.Infoln(plugin.Name())
			plugins = append(plugins, plugin)
		} else {
			klog.Infoln(i, " : no such filter algorithm")
		}
	}

	/*
		for _, cluster := range *clusterInfoList {
			score = 0
			klog.Infoln(">>", cluster.ClusterName)
			if cluster.IsFiltered {
				klog.Infoln(cluster.ClusterName, "is already filtered")
			} else {
				for _, plugin := range plugins {
					klog.Infoln("[plugin]", plugin.Name())
					score = plugin.Score(pod, status, cluster)
					klog.Infoln(score)
					result[cluster.ClusterName] += score
				}
			}
			(*cluster).ClusterScore = int32(result[cluster.ClusterName])
		}
	*/
	for _, plugin := range plugins {
		score = 0
		tmp := new(util.TmpEachScore)
		f.pluginScoreList[plugin.Name()] = *tmp
		klog.Infoln(">>", plugin.Name())
		for _, cluster := range *clusterInfoList {
			klog.Infoln(">>", cluster.ClusterName)
			if !cluster.IsFiltered {
				score = plugin.Score(pod, status, cluster)
				klog.Infoln(score)
				tmp.Total += score
				tmp.ScoreList[cluster.ClusterName] = score
			}
			plugin.Normalize(tmp, clusterInfoList)
		}
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
