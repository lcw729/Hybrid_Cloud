package main

import "Hybrid_Cloud/hcp-scheduler/pkg/scheduler"

func main() {
	sched := scheduler.NewScheduler()
	sched.TestScheduling()
}
