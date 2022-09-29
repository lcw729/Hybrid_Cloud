package predicates

import (
	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

type TaintToleration struct {
	Temp           []int64
	IntolerableMap map[string]int64
}

func (pl *TaintToleration) Name() string {
	return plugins.TaintToleration
}

// ComputeTaintTolerationPriorityMap prepares the priority list for all the nodes based on the number of intolerable taints on the node
func (pl *TaintToleration) Filter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	for _, node := range clusterInfo.Nodes {
		if node.Node == nil {
			continue
		}

		filterPredicate := func(t *v1.Taint) bool {
			// PodToleratesNodeTaints is only interested in NoSchedule and NoExecute taints.
			return t.Effect == v1.TaintEffectNoSchedule || t.Effect == v1.TaintEffectNoExecute
		}

		taint, isUntolerated := FindMatchingUntoleratedTaint(node.Node.Spec.Taints, pod.Spec.Tolerations, filterPredicate)
		if !isUntolerated {
			return false
		} else {
			klog.Infof("node(s) had taint {%s: %s}, that the pod didn't tolerate\n", taint.Key, taint.Value)
		}
	}
	return true
}

type taintsFilterFunc func(*v1.Taint) bool

// FindMatchingUntoleratedTaint checks if the given tolerations tolerates
// all the filtered taints, and returns the first taint without a toleration
func FindMatchingUntoleratedTaint(taints []v1.Taint, tolerations []v1.Toleration, inclusionFilter taintsFilterFunc) (v1.Taint, bool) {
	filteredTaints := getFilteredTaints(taints, inclusionFilter)
	for _, taint := range filteredTaints {
		if !TolerationsTolerateTaint(tolerations, &taint) {
			return taint, true
		}
	}
	return v1.Taint{}, false
}

// TolerationsTolerateTaint checks if taint is tolerated by any of the tolerations.
func TolerationsTolerateTaint(tolerations []v1.Toleration, taint *v1.Taint) bool {
	for i := range tolerations {
		if tolerations[i].ToleratesTaint(taint) {
			return true
		}
	}
	return false
}

// getFilteredTaints returns a list of taints satisfying the filter predicate
func getFilteredTaints(taints []v1.Taint, inclusionFilter taintsFilterFunc) []v1.Taint {
	if inclusionFilter == nil {
		return taints
	}
	filteredTaints := []v1.Taint{}
	for _, taint := range taints {
		if !inclusionFilter(&taint) {
			continue
		}
		filteredTaints = append(filteredTaints, taint)
	}
	return filteredTaints
}

/*
func (pl *TaintToleration) Normalize() {
	max := pl.Temp[0]
	for i, j := range pl.IntolerableMap {
		if j == 0 {
			pl.IntolerableMap[i] = scoretable.MaxNodeScore
		} else {
			pl.IntolerableMap[i] = int64(100 * ((float32(max) - float32(pl.IntolerableMap[i])) / float32(max)))
		}
	}
}

func (pl *TaintToleration) Sort() {
	sort.Slice(pl.Temp, func(i, j int) bool {
		return pl.Temp[i] > pl.Temp[j]
	})
}

// CountIntolerableTaintsPreferNoSchedule gives the count of intolerable taints of a pod with effect PreferNoSchedule
func countIntolerableTaintsPreferNoSchedule(taints []v1.Taint, tolerations []v1.Toleration) (intolerableTaints int) {
	for _, taint := range taints {
		// check only on taints that have effect PreferNoSchedule
		if taint.Effect != v1.TaintEffectPreferNoSchedule {
			continue
		}

		// 해당 toleration이 존재하지 않으면 intolerableTaints + 1
		if !TolerationsTolerateTaint(tolerations, &taint) {
			intolerableTaints++
		}
	}
	return
}

// Taint에 해당하는 toleration이 존재하는지 확인
// false - 존재하지 않음
// true - 존재함
func TolerationsTolerateTaint(tolerations []v1.Toleration, taint *v1.Taint) bool {
	for i := range tolerations {
		if tolerations[i].ToleratesTaint(taint) {
			return true
		}
	}
	return false
}

// getAllTolerationEffectPreferNoSchedule gets the list of all Tolerations with Effect PreferNoSchedule or with no effect.
func getAllTolerationPreferNoSchedule(tolerations []v1.Toleration) (tolerationList []v1.Toleration) {
	for _, toleration := range tolerations {
		// Empty effect means all effects which includes PreferNoSchedule, so we need to collect it as well.
		if len(toleration.Effect) == 0 || toleration.Effect == v1.TaintEffectPreferNoSchedule {
			tolerationList = append(tolerationList, toleration)
		}
	}
	return
}
*/
