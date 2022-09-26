package predicates

import (
	"github.com/KETI-Hybrid/hcp-pkg/kube-resource/kubefed"

	"github.com/KETI-Hybrid/hcp-scheduler-v1/src/framework/plugins"
	"github.com/KETI-Hybrid/hcp-scheduler-v1/src/resourceinfo"

	v1 "k8s.io/api/core/v1"
)

type JoinCheck struct{}

func (pl *JoinCheck) Name() string {
	return plugins.JoinCheck
}

func (pl *JoinCheck) Filter(pod *v1.Pod, status *resourceinfo.CycleStatus, clusterInfo *resourceinfo.ClusterInfo) bool {
	return !kubefed.IsKubeFedCluster(clusterInfo.ClusterName)
}
