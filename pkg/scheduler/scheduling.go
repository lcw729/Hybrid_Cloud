package scheduler

import (
	policy "Hybrid_Cloud/hcp-resource/hcppolicy"
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/priorities"
	"Hybrid_Cloud/hcp-scheduler/pkg/internal/scoretable"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Scheduler watches for new unscheduled pods. It attempts to find
// nodes that they fit on and writes bindings back to the api server.
type Scheduler struct {
	SchedulingResource v1.Pod
	ClusterClients     map[string]*kubernetes.Clientset
	ClusterInfoList    *resourceinfo.ClusterInfoList
	ScoreTable         *map[string]float32
	SchdPolicy         string
}

var AlgorithmMap = map[string]func(*v1.Pod, *v1.Node) int32{

	"Affinity":        priorities.NodeAffinity,
	"TaintToleration": priorities.TaintToleration,
}

/*
func ListTargetClusters() []string {
	cm, _ := clusterManager.NewClusterManager()
	kubefed_clusters := cm.Cluster_list.Items
	var target_clusters []string
	for _, i := range kubefed_clusters {
		target_clusters = append(target_clusters, i.ObjectMeta.Name)
	}
	return target_clusters
}
*/

func NewScheduler() *Scheduler {
	cm, _ := clusterManager.NewClusterManager()

	schd := Scheduler{
		ClusterClients:  cm.Cluster_kubeClients,
		ClusterInfoList: resourceinfo.NewClusterInfoList(),
		// ClusterList:    ListTargetClusters(),
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

func (sched *Scheduler) Scheduling(deployment *v1alpha1.HCPDeployment) []v1alpha1.Target {

	schedule_type := sched.SchdPolicy
	replicas := *deployment.Spec.RealDeploymentSpec.Replicas
	var scheduling_result []v1alpha1.Target
	var cnt int32 = 0

	// set schedulingResource
	fake_pod := newPodFromHCPDeployment(deployment)
	sched.SchedulingResource = *fake_pod

	// set scoretable
	sched.ScoreTable = scoretable.NewScoreTable(len(*sched.ClusterInfoList))

	for i := 0; i < int(replicas); i++ {
		sched.Scoring(schedule_type)
		score_table := sched.ScoreTable
		if score_table != nil {
			target := scoretable.SortScore(*score_table)[0].Cluster
			if target != "" {
				if updateSchedulingResult(&scheduling_result, target) {
					cnt += 1
					fmt.Printf("[Scheduling] %d/%d pod / TargetCluster : %s\n", i+1, replicas, target)
					fmt.Println()
				}
			} else {
				fmt.Println("ERROR: No cluster to be scheduled")
				fmt.Println("Scheduling failed")
				break
			}
		} else {
			fmt.Println("Scheduling failed")
			return nil
		}
	}

	if cnt == replicas {
		fmt.Println("Scheduling succeeded")
		printSchedulingResult(scheduling_result)
		return scheduling_result
	} else {
		fmt.Println("Scheduling failed")
		return nil
	}
}

func printSchedulingResult(targets []v1alpha1.Target) {
	fmt.Println("========scheduling result========")
	for _, i := range targets {
		fmt.Println("target cluster :", i.Cluster)
		fmt.Println("replicas       :", *i.Replicas)
		fmt.Println()
	}
}

func newPodFromHCPDeployment(deployment *v1alpha1.HCPDeployment) *v1.Pod {

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        deployment.GetObjectMeta().GetName() + "-pod",
			Annotations: deployment.Annotations,
			Labels:      deployment.Labels,
		},
		Spec: deployment.Spec.RealDeploymentSpec.Template.Spec,
	}

	return pod
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

func (sched *Scheduler) Scoring(algorithm string) {

	var pod = &sched.SchedulingResource
	clusterinfoList := sched.ClusterInfoList
	var score int32

	for _, clusterinfo := range *clusterinfoList {
		fmt.Println("==>", clusterinfo.ClusterName)
		score = 0
		node_list := (*clusterinfo).Nodes

		for _, node := range node_list {
			node_score := AlgorithmMap[algorithm](pod, node.Node)
			if node_score == -1 {
				fmt.Println("fail to scoring node")
				return
			} else {
				fmt.Println(node.NodeName, "score :", node_score)
				score += node_score
			}
		}
		fmt.Println("*", clusterinfo.ClusterName, "total score :", score)
		(*sched.ScoreTable)[clusterinfo.ClusterName] = float32(score)
	}
}

func (sched *Scheduler) scheduleOne(ctx context.Context) {

}
