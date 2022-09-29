/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package priorities

import (
	"fmt"

	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/resourceinfo"
	"hcp-scheduler/src/util"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/klog"
)

type NodeAffinity struct{}

func (pl *NodeAffinity) Name() string {
	return plugins.NodeAffinity
}

// CalculateNodeAffinityPriorityMap prioritizes nodes according to node affinity scheduling preferences
// indicated in PreferredDuringSchedulingIgnoredDuringExecution. Each time a node match a preferredSchedulingTerm,
// it will a get an add of preferredSchedulingTerm.Weight. Thus, the more preferredSchedulingTerms
// the node satisfies and the more the preferredSchedulingTerm that is satisfied weights, the higher
// score the node gets.
func (pl *NodeAffinity) Score(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) int64 {
	var count int32
	for _, node := range clusterInfo.Nodes {
		// default is the podspec.
		affinity := pod.Spec.Affinity
		// A nil element of PreferredDuringSchedulingIgnoredDuringExecution matches no objects.
		// An element of PreferredDuringSchedulingIgnoredDuringExecution that refers to an
		// empty PreferredSchedulingTerm matches all objects.
		if affinity != nil && affinity.NodeAffinity != nil && affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
			// Match PreferredDuringSchedulingIgnoredDuringExecution term by term.
			for i := range affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution {

				preferredSchedulingTerm := &affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[i]
				if preferredSchedulingTerm.Weight == 0 {
					continue
				}

				// TODO: Avoid computing it for all nodes if this becomes a performance problem.
				nodeSelector, err := NodeSelectorRequirementsAsSelector(preferredSchedulingTerm.Preference.MatchExpressions)
				if err != nil {
					klog.Error(err)
					return -1
				}

				if nodeSelector.Matches(labels.Set((*node.Node).Labels)) {
					count += preferredSchedulingTerm.Weight
				}
			}
		}
	}
	return int64(count)
}

func (pl *NodeAffinity) Normalize(tmpEachScore *util.TmpEachScore, clusterInfoList *resourceinfo.ClusterInfoList) {
	for _, cluster := range *clusterInfoList {
		klog.Infoln(">>", cluster.ClusterName)
		if !cluster.IsFiltered {
			tmpEachScore.ScoreList[cluster.ClusterName] /= tmpEachScore.Total
			fmt.Println(tmpEachScore.ScoreList[cluster.ClusterName])
			(*cluster).ClusterScore += int32(tmpEachScore.ScoreList[cluster.ClusterName])
		}
	}
}

// NodeSelectorRequirementsAsSelector converts the []NodeSelectorRequirement api type into a struct that implements
// labels.Selector.
func NodeSelectorRequirementsAsSelector(nsm []v1.NodeSelectorRequirement) (labels.Selector, error) {
	if len(nsm) == 0 {
		return labels.Nothing(), nil
	}
	selector := labels.NewSelector()
	for _, expr := range nsm {
		var op selection.Operator
		switch expr.Operator {
		case v1.NodeSelectorOpIn:
			op = selection.In
		case v1.NodeSelectorOpNotIn:
			op = selection.NotIn
		case v1.NodeSelectorOpExists:
			op = selection.Exists
		case v1.NodeSelectorOpDoesNotExist:
			op = selection.DoesNotExist
		case v1.NodeSelectorOpGt:
			op = selection.GreaterThan
		case v1.NodeSelectorOpLt:
			op = selection.LessThan
		default:
			return nil, fmt.Errorf("%q is not a valid node selector operator", expr.Operator)
		}
		r, err := labels.NewRequirement(expr.Key, op, expr.Values)
		if err != nil {
			return nil, err
		}
		selector = selector.Add(*r)
	}
	return selector, nil
}
