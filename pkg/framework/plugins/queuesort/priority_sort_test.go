/*
Copyright 2020 The Kubernetes Authors.

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

package queuesort

import (
	"Hybrid_Cloud/kube-resource/pod"
	"fmt"
	"testing"

	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func TestLess(t *testing.T) {
	prioritySort := &PrioritySort{}
	high_pod, _ := pod.GetPod("kube-master", "high-nginx", "default")
	higher_pod, _ := pod.GetPod("kube-master", "higher-nginx", "default")
	// var lowPriority, highPriority = int32(10), int32(100)
	//t1 := time.Now()
	//t2 := t1.Add(time.Second)
	type TestQueueSort struct {
		name     string
		p1       *framework.QueuedPodInfo
		p2       *framework.QueuedPodInfo
		expected bool
	}

	tt := TestQueueSort{
		name: "p1.priority less than p2.priority",
		p1: &framework.QueuedPodInfo{
			PodInfo: framework.NewPodInfo(high_pod),
		},
		p2: &framework.QueuedPodInfo{
			PodInfo: framework.NewPodInfo(higher_pod),
		},
		expected: false,
	}

	t.Run(tt.name, func(t *testing.T) {
		if got := prioritySort.Less(tt.p1, tt.p2); got != tt.expected {
			t.Errorf("expected %v, got %v", tt.expected, got)
		} else {
			fmt.Println(tt.p1.PodInfo.Pod.ObjectMeta.Name)
			fmt.Println(tt.p2.PodInfo.Pod.ObjectMeta.Name)
		}
	})
}
