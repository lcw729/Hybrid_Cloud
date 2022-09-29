package priorities

import (
	"fmt"
	"math"

	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/internal/scoretable"
	"hcp-scheduler/src/resourceinfo"
	"hcp-scheduler/src/util"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type BalanceAllocation struct{}

func (pl *BalanceAllocation) Name() string {
	return plugins.BalanceAllocation
}

func (pl *BalanceAllocation) Score(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) int64 {

	var score int64 = 0
	for _, node := range clusterInfo.Nodes {

		requested := util.CreateResourceToValueMapPO(pod)
		klog.Infoln("requested CPU:", requested[v1.ResourceCPU], "requested Memory:", requested[v1.ResourceMemory])
		allocable := util.CreateResourceToValueMapNode(node.Node)
		klog.Infoln("allocable CPU:", allocable[v1.ResourceCPU], "allocable Memory:", allocable[v1.ResourceMemory])

		var includeVolumes bool
		if len(node.Node.Status.VolumesAttached) > 0 {
			includeVolumes = true
		}
		requestedVolumes := len(pod.Spec.Volumes)
		allocatableVolumes := len(node.Node.Status.VolumesAttached) - len(node.Node.Status.VolumesInUse)

		score += balancedResourceScorer(requested, allocable, includeVolumes, requestedVolumes, allocatableVolumes)

	}

	return score
}

func (pl *BalanceAllocation) Normalize(tmpEachScore *util.TmpEachScore, clusterInfoList *resourceinfo.ClusterInfoList) {
	for _, cluster := range *clusterInfoList {
		klog.Infoln(">>", cluster.ClusterName)
		if !cluster.IsFiltered {
			tmpEachScore.ScoreList[cluster.ClusterName] /= tmpEachScore.Total
			fmt.Println(tmpEachScore.ScoreList[cluster.ClusterName])
			(*cluster).ClusterScore += int32(tmpEachScore.ScoreList[cluster.ClusterName])
		}
	}
}

// todo: use resource weights in the scorer function
func balancedResourceScorer(requested, allocable util.ResourceToValueMap, includeVolumes bool, requestedVolumes int, allocatableVolumes int) int64 {

	// capacity 대비 request의 비율
	cpuFraction := fractionOfCapacity(requested[v1.ResourceCPU], allocable[v1.ResourceCPU])
	memoryFraction := fractionOfCapacity(requested[v1.ResourceMemory], allocable[v1.ResourceMemory])
	if cpuFraction >= 1 || memoryFraction >= 1 {
		// if requested >= capacity, the corresponding host should never be preferred.
		return 0
	}

	// volume 요청한 경우
	if includeVolumes && allocatableVolumes > 0 {

		volumeFraction := float64(requestedVolumes) / float64(allocatableVolumes)

		if volumeFraction >= 1 {
			return 0
		} else {
			mean := (cpuFraction + memoryFraction + volumeFraction) / float64(3)
			variance := float64((((cpuFraction - mean) * (cpuFraction - mean)) + ((memoryFraction - mean) * (memoryFraction - mean)) + ((volumeFraction - mean) * (volumeFraction - mean))) / float64(3))
			return int64((1 - variance) * float64(scoretable.MaxNodeScore))
		}
	}

	diff := math.Abs(cpuFraction - memoryFraction)
	return int64((1 - diff) * float64(scoretable.MaxNodeScore))
}

func fractionOfCapacity(requested, capacity int64) float64 {
	if capacity == 0 {
		return 1
	}
	return float64(requested) / float64(capacity)
}
