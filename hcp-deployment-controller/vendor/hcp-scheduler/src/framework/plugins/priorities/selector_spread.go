package priorities

import (
	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

// SelectorSpread is a plugin that calculates selector spread priority.
type SelectorSpread struct{}

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = plugins.SelectorSpread
	// preScoreStateKey is the key in CycleState to SelectorSpread pre-computed data for Scoring.
	preScoreStateKey = "PreScore" + Name
)

// preScoreState computed at PreScore and used at Score.
type preScoreState struct {
	selector labels.Selector
}

// Clone implements the mandatory Clone interface. We don't really copy the data since
// there is no need for that.
func (s *preScoreState) Clone() resourceinfo.StateData {
	return s
}

// skipSelectorSpread returns true if the pod's TopologySpreadConstraints are specified.
// Note that this doesn't take into account default constraints defined for
// the PodTopologySpread plugin.
func skipSelectorSpread(pod *v1.Pod) bool {
	return len(pod.Spec.TopologySpreadConstraints) != 0
}

// Score invoked at the Score extension point.
// The "score" returned in this function is the matching number of pods on the `nodeName`,
// it is normalized later.
func (pl *SelectorSpread) Score(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) int64 {
	var count = 0
	for _, nodeInfo := range clusterInfo.Nodes {
		if skipSelectorSpread(pod) {
			return 0
		}

		c, err := status.Read(preScoreStateKey)
		if err != nil {
			return 0
		}

		s, ok := c.(*preScoreState)
		if !ok {
			return 0
		}

		count += countMatchingPods(pod.Namespace, s.selector, nodeInfo)
	}
	return int64(count)
}

// countMatchingPods counts pods based on namespace and matching all selectors
func countMatchingPods(namespace string, selector labels.Selector, nodeInfo *resourceinfo.NodeInfo) int {
	if len(nodeInfo.Pods) == 0 || selector.Empty() {
		return 0
	}
	count := 0
	for _, p := range nodeInfo.Pods {
		// Ignore pods being deleted for spreading purposes
		// Similar to how it is done for SelectorSpreadPriority
		if namespace == p.Pod.Namespace && p.Pod.DeletionTimestamp == nil {
			if selector.Matches(labels.Set(p.Pod.Labels)) {
				count++
			}
		}
	}
	return count
}
