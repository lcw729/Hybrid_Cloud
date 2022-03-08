package pkg

import (
	"context"

	"k8s.io/client-go/kubernetes"
)

type Scheduler struct {
	ClusterClients map[string]*kubernetes.Clientset
	SchdPolicy     string
}

func New() {

}

// scheduleOne does the entire scheduling workflow for a single pod. It is serialized on the scheduling algorithm's host fitting.
func (sched *Scheduler) scheduleOne(ctx context.Context) {

}

func (sched *Scheduler) Scheduling() {

}
