package predicates

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"

	v1 "k8s.io/api/core/v1"
)

// Filter invoked at the filter extension point.
func CheckNodePorts(pod *v1.Pod, nodeInfo *resourceinfo.NodeInfo) bool {
	fits := fitsPorts(getContainerPorts(pod), nodeInfo)
	if !fits {
		return false
	} else {
		return true
	}
}

// getContainerPorts returns the used host ports of Pods: if 'port' was used, a 'port:true' pair
// will be in the result; but it does not resolve port conflict.
func getContainerPorts(pods ...*v1.Pod) []*v1.ContainerPort {
	ports := []*v1.ContainerPort{}
	for _, pod := range pods {
		for j := range pod.Spec.Containers {
			container := &pod.Spec.Containers[j]
			for k := range container.Ports {
				ports = append(ports, &container.Ports[k])
			}
		}
	}
	return ports
}

func fitsPorts(wantPorts []*v1.ContainerPort, nodeInfo *resourceinfo.NodeInfo) bool {
	// try to see whether existingPorts and wantPorts will conflict or not
	existingPorts := nodeInfo.UsedPorts
	for _, cp := range wantPorts {
		if existingPorts.CheckConflict(cp.HostIP, string(cp.Protocol), cp.HostPort) {
			return false
		}
	}
	return true
}
