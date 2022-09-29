package predicates

import (
	"hcp-pkg/kube-resource/kubefed"

	"hcp-scheduler/src/framework/plugins"
	"hcp-scheduler/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
)

type JoinCheck struct{}

func (pl *JoinCheck) Name() string {
	return plugins.JoinCheck
}

func (pl *JoinCheck) Filter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	return !kubefed.IsKubeFedCluster(clusterInfo.ClusterName)
}
