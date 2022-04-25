package util

import (
	v1 "k8s.io/api/core/v1"
)

// ResourceToValueMap contains resource name and score.
type ResourceToValueMap map[v1.ResourceName]int64

func CreateResourceToValueMapPO(pod *v1.Pod) ResourceToValueMap {

	var resourceMap = make(ResourceToValueMap)
	resourceMap[v1.ResourceCPU] = getMilliCPU(pod)
	resourceMap[v1.ResourceMemory] = getMilliMemory(pod)
	resourceMap[v1.ResourceEphemeralStorage] = getMilliEphemeralStorage(pod)

	return resourceMap
}

func CreateResourceToValueMapNode(node *v1.Node) ResourceToValueMap {

	var resourceMap = make(ResourceToValueMap)
	resourceMap[v1.ResourceCPU] = getAllocableCPU(node)
	resourceMap[v1.ResourceMemory] = getAllocableMemory(node)

	return resourceMap
}

func getMilliCPU(pod *v1.Pod) int64 {
	var cpu int64 = 0
	containers := pod.Spec.Containers
	for _, c := range containers {
		cpu += c.Resources.Requests.Cpu().MilliValue()
	}

	return cpu
}

func getMilliCPULimit(pod *v1.Pod) int64 {
	var cpu_limit int64 = 0
	containers := pod.Spec.Containers
	for _, c := range containers {
		cpu_limit += c.Resources.Limits.Cpu().MilliValue()
	}

	return cpu_limit
}

func getMilliMemory(pod *v1.Pod) int64 {
	var mem int64 = 0
	containers := pod.Spec.Containers
	for _, c := range containers {
		mem += c.Resources.Limits.Memory().MilliValue()
	}

	return mem
}

func getMilliMemoryLimit(pod *v1.Pod) int64 {
	var mem_limit int64 = 0
	containers := pod.Spec.Containers
	for _, c := range containers {
		mem_limit += c.Resources.Limits.Memory().MilliValue()
	}

	return mem_limit
}

func getMilliEphemeralStorage(pod *v1.Pod) int64 {
	var total int64 = 0
	containers := pod.Spec.Containers
	for _, c := range containers {
		total += c.Resources.Requests.StorageEphemeral().MilliValue()
	}

	return total
}

func getAllocableCPU(node *v1.Node) int64 {
	cpu := node.Status.Allocatable.Cpu().MilliValue()
	return cpu
}

func getAllocableMemory(node *v1.Node) int64 {
	mem := node.Status.Allocatable.Memory().MilliValue()
	return mem
}
