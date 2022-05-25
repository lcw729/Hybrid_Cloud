package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/framework/plugins"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/kube-resource/kubefed"

	v1 "k8s.io/api/core/v1"
)

type JoinCheck struct{}

func (pl *JoinCheck) Name() string {
	return plugins.JoinCheck
}

func (pl *JoinCheck) Filter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	return kubefed.Iskubefedcluster(clusterInfo.ClusterName)
}
