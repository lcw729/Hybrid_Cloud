package priorities

import (
	"fmt"
	"strings"

	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/internal/scoretable"
	"hcp-scheduler/src/resourceinfo"
	"hcp-scheduler/src/util"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

// The two thresholds are used as bounds for the image score range. They correspond to a reasonable size range for
// container images compressed and stored in registries; 90%ile of images on dockerhub drops into this range.
const (
	mb           int64 = 1024 * 1024
	minThreshold int64 = 23 * mb
	maxThreshold int64 = 1000 * mb
)

type ImageLocality struct{}

func (pl *ImageLocality) Name() string {
	return plugins.ImageLocality
}

// ImageLocalityPriorityMap is a priority function that favors nodes that already have requested pod container's images.
// It will detect whether the requested images are present on a node, and then calculate a score ranging from 0 to 10
// based on the total size of those images.
// - If none of the images are present, this node will be given the lowest priority.
// - If some of the images are present on a node, the larger their sizes' sum, the higher the node's priority.
func (pl *ImageLocality) Score(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) int64 {
	var score int64 = 0
	for _, node := range clusterInfo.Nodes {
		if node == nil {
			return 0
		}
		score += calculatePriority(sumImageScores(node, pod.Spec.Containers, status.TotalNumNodes))
	}
	return int64(score)
}

func (pl *ImageLocality) Normalize(tmpEachScore *util.TmpEachScore, clusterInfoList *resourceinfo.ClusterInfoList) {
	for _, cluster := range *clusterInfoList {
		klog.Infoln(">>", cluster.ClusterName)
		if !cluster.IsFiltered {
			tmpEachScore.ScoreList[cluster.ClusterName] /= tmpEachScore.Total
			fmt.Println(tmpEachScore.ScoreList[cluster.ClusterName])
			(*cluster).ClusterScore += int32(tmpEachScore.ScoreList[cluster.ClusterName])
		}
	}
}

// calculatePriority returns the priority of a node. Given the sumScores of requested images on the node, the node's
// priority is obtained by scaling the maximum priority value with a ratio proportional to the sumScores.
func calculatePriority(sumScores int64) int64 {
	if sumScores < minThreshold {
		sumScores = minThreshold
	} else if sumScores > maxThreshold {
		sumScores = maxThreshold
	}

	return int64(scoretable.MaxNodeScore) * (sumScores - minThreshold) / (maxThreshold - minThreshold)
}

// sumImageScores returns the sum of image scores of all the containers that are already on the node.
// Each image receives a raw score of its size, scaled by scaledImageScore. The raw scores are later used to calculate
// the final score. Note that the init containers are not considered for it's rare for users to deploy huge init containers.
func sumImageScores(nodeInfo *resourceinfo.NodeInfo, containers []v1.Container, totalNumNodes int) int64 {
	var sum int64
	imageStates := nodeInfo.GetImageStates()

	for _, container := range containers {
		if state, ok := imageStates[normalizedImageName(container.Image)]; ok {
			sum += scaledImageScore(state, totalNumNodes)
		}
	}

	return sum
}

// scaledImageScore returns an adaptively scaled score for the given state of an image.
// The size of the image is used as the base score, scaled by a factor which considers how much nodes the image has "spread" to.
// This heuristic aims to mitigate the undesirable "node heating problem", i.e., pods get assigned to the same or
// a few nodes due to image locality.
func scaledImageScore(imageState *resourceinfo.ImageStateSummary, totalNumNodes int) int64 {
	spread := float64(imageState.NumNodes) / float64(totalNumNodes)
	return int64(float64(imageState.Size) * spread)
}

// normalizedImageName returns the CRI compliant name for a given image.
// TODO: cover the corner cases of missed matches, e.g,
// 1. Using Docker as runtime and docker.io/library/test:tag in pod spec, but only test:tag will present in node status
// 2. Using the implicit registry, i.e., test:tag or library/test:tag in pod spec but only docker.io/library/test:tag
// in node status; note that if users consistently use one registry format, this should not happen.
func normalizedImageName(name string) string {
	if strings.LastIndex(name, ":") <= strings.LastIndex(name, "/") {
		name = name + ":latest"
	}

	return name
}
