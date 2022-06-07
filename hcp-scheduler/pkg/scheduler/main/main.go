package main

import (
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"fmt"
)

func main() {
	clusterInfoList := resourceinfo.NewClusterInfoList()
	clusterInfoMap := resourceinfo.CreateClusterInfoMap(clusterInfoList)
	fmt.Println(clusterInfoMap["eks-cluster"].Nodes[0].Pods[0].Pod.Name)
	// sched := scheduler.NewScheduler()
	// sched.Scheduling(nil)
}
