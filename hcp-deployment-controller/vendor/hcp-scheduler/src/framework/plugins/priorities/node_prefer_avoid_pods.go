package priorities

import (
	"encoding/json"

	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/internal/scoretable"
	"hcp-scheduler/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type NodePreferAvoidPods struct{}

func (pl *NodePreferAvoidPods) Name() string {
	return plugins.NodePreferAvoidPods
}
func (pl *NodePreferAvoidPods) Score(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) int64 {
	var score int64 = 0
	for _, node := range clusterInfo.Nodes {

		controllerRef := metav1.GetControllerOf(pod)
		if controllerRef != nil {
			// Ignore pods that are owned by other controller than ReplicationController
			// or ReplicaSet.
			if controllerRef.Kind != "ReplicationController" && controllerRef.Kind != "ReplicaSet" {
				controllerRef = nil
			}
		}
		if controllerRef == nil {
			score += scoretable.MaxNodeScore
		}

		avoids, err := GetAvoidPodsFromNodeAnnotations(node.Node.Annotations)
		if err != nil {
			// If we cannot get annotation, assume it's schedulable there.
			score += scoretable.MaxNodeScore
		}
		for i := range avoids.PreferAvoidPods {
			avoid := &avoids.PreferAvoidPods[i]
			if avoid.PodSignature.PodController.Kind == controllerRef.Kind && avoid.PodSignature.PodController.UID == controllerRef.UID {
				break
			}
		}

		score += scoretable.MaxNodeScore
	}

	return score
}

// GetAvoidPodsFromNodeAnnotations scans the list of annotations and
// returns the pods that needs to be avoided for this node from scheduling
func GetAvoidPodsFromNodeAnnotations(annotations map[string]string) (v1.AvoidPods, error) {
	var avoidPods v1.AvoidPods
	if len(annotations) > 0 && annotations[v1.PreferAvoidPodsAnnotationKey] != "" {
		err := json.Unmarshal([]byte(annotations[v1.PreferAvoidPodsAnnotationKey]), &avoidPods)
		if err != nil {
			return avoidPods, err
		}
	}
	return avoidPods, nil
}
