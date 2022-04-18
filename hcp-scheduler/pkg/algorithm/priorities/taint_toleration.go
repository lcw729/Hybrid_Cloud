package priorities

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/internal/scoretable"
	"fmt"

	v1 "k8s.io/api/core/v1"
)

// CountIntolerableTaintsPreferNoSchedule gives the count of intolerable taints of a pod with effect PreferNoSchedule
func countInTolerableTaintsPreferNoSchedule(taints []v1.Taint, tolerations []v1.Toleration) (intolerableTaints int) {
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
	return intolerableTaints
}

// Taint에 해당하는 toleration이 존재하는지 확인
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

// ComputeTaintTolerationPriorityMap prepares the priority list for all the nodes based on the number of intolerable taints on the node
func TaintToleration(pod *v1.Pod, node *v1.Node) int32 {

	if node == nil {
		return -1
	}
	// To hold all the tolerations with Effect PreferNoSchedule
	var tolerationsPreferNoSchedule = getAllTolerationPreferNoSchedule(pod.Spec.Tolerations)
	score := int32(countInTolerableTaintsPreferNoSchedule(node.Spec.Taints, tolerationsPreferNoSchedule))

	// 모두 intolerable인 경우, MinNodeScore
	if score > 0 {
		taints_len := int32(len(node.Spec.Taints))
		if score == taints_len {
			score = scoretable.MinNodeScore
		} else {
			fmt.Println("here")
			score = int32(100 * float64(score) / float64(taints_len))
		}
	} else if score == 0 {
		// intolerable인 경우, 개수에 따라 차등 점수
		score = scoretable.MaxNodeScore
	}

	return score
}
