package scheduler

import (
	policy "Hybrid_Cloud/hcp-resource/hcppolicy"
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"

	"k8s.io/client-go/kubernetes"
)

// Scheduler watches for new unscheduled pods. It attempts to find
// nodes that they fit on and writes bindings back to the api server.
type Scheduler struct {
	ClusterClients map[string]*kubernetes.Clientset
	ClusterInfo    []*resourceinfo.ClusterInfo
	ClusterList    []string
	SchdPolicy     string
}

func ListTargetClusters() []string {
	cm, _ := clusterManager.NewClusterManager()
	kubefed_clusters := cm.Cluster_list.Items
	var target_clusters []string
	for _, i := range kubefed_clusters {
		target_clusters = append(target_clusters, i.ObjectMeta.Name)
	}
	return target_clusters
}

func NewScheduler() *Scheduler {
	cm, _ := clusterManager.NewClusterManager()

	schd := Scheduler{
		ClusterClients: cm.Cluster_kubeClients,
		ClusterInfo:    resourceinfo.NewClusterInfoList(),
		ClusterList:    ListTargetClusters(),
	}

	// HCPPolicy 최적 배치 알고리즘 정책 읽어오기
	algorithm, err := policy.GetAlgorithm()
	if err == nil {
		schd.SchdPolicy = algorithm
	} else {
		schd.SchdPolicy = "DEFAULT_SCHEDPOLICY"
	}

	return &schd
}

func (s *Scheduler) Scheduling(deployment *v1alpha1.HCPDeployment) []v1alpha1.Target {

	schedule_type := s.SchdPolicy
	replicas := *deployment.Spec.RealDeploymentSpec.Replicas
	var scheduling_result []v1alpha1.Target
	var cnt int32 = 0

	for i := 0; i < int(replicas); i++ {
		switch schedule_type {
		case "Affinity":
			target := algorithm.Affinity(&s.ClusterList)
			if target != "" {
				if updateSchedulingResult(&scheduling_result, target) {
					cnt += 1
					fmt.Printf("[Scheduling] %d/%d pod / TargetCluster : %s\n", i+1, replicas, target)
				}
			} else {
				fmt.Println("ERROR: No cluster to be scheduled")
				fmt.Println("Scheduling failed")
				return nil
			}
		case "DRF":
		}
	}

	if cnt == replicas {
		fmt.Println("Scheduling succeeded")
		fmt.Println("[Scheduling Result] =====> ", scheduling_result[0].Cluster, *scheduling_result[0].Replicas)
		return scheduling_result
	} else {
		fmt.Println("Scheduling failed")
		return nil
	}
}

func updateSchedulingResult(scheduling_result *[]v1alpha1.Target, cluster string) bool {

	for i, target := range *scheduling_result {
		// 이미 target cluster 목록에 cluster가 있는 경우
		if target.Cluster == cluster {
			// replicas 개수 증가
			temp := *target.Replicas
			temp += 1
			target.Replicas = &temp
			(*scheduling_result)[i] = target
			return true
		}
	}

	// target cluster 목록에 cluster가 없는 경우

	// replicas 개수 1로 설정
	var new_target v1alpha1.Target
	new_target.Cluster = cluster
	var one int32 = 1
	new_target.Replicas = &one
	*scheduling_result = append((*scheduling_result), new_target)
	return true
}

func (sched *Scheduler) scheduleOne(ctx context.Context) {

}
