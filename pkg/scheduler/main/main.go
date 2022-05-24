package main

import "Hybrid_Cloud/hcp-scheduler/pkg/scheduler"

func main() {
	// pod := &v1.Pod{
	// 	Spec: v1.PodSpec{
	// 		Volumes: []v1.Volume{
	// 			{
	// 				VolumeSource: v1.VolumeSource{
	// 					HostPath: &v1.HostPathVolumeSource{
	// 						Path: "/test",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// podVols := pod.Spec.Volumes
	// for _, podVol := range podVols {
	// 	fmt.Println(podVol.AWSElasticBlockStore)
	// }
	sched := scheduler.NewScheduler()
	sched.Scheduling(nil)
}
