package scheduler

import (
	policy "Hybrid_Cloud/hcp-resource/hcppolicy"
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/predicates"
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/priorities"
	"Hybrid_Cloud/hcp-scheduler/pkg/algorithm/test"
	"Hybrid_Cloud/hcp-scheduler/pkg/internal/scoretable"
	"Hybrid_Cloud/hcp-scheduler/pkg/resourceinfo"
	"Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"
	"sort"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Scheduler watches for new unscheduled pods. It attempts to find
// nodes that they fit on and writes bindings back to the api server.
type Scheduler struct {
	SchedulingResource *v1.Pod
	ClusterClients     map[string]*kubernetes.Clientset
	ClusterInfoList    resourceinfo.ClusterInfoList
	//NodeInfoMap        map[string]*resourceinfo.NodeInfo
	ClusterInfoMap   map[string]*resourceinfo.ClusterInfo
	SchedulingResult []v1alpha1.Target
	SchedPolicy      string
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

	//	clusterInfoList := *resourceinfo.NewClusterInfoList()
	clusterInfoList := &resourceinfo.ClusterInfoList{}
	test.CreateTestClusterNodeResourcesBalancedAllocation(clusterInfoList)
	clusterInfoMap := resourceinfo.CreateClusterInfoMap(clusterInfoList)
	// HCPPolicy 최적 배치 알고리즘 정책 읽어오기
	algorithm, err := policy.GetAlgorithm()
	var schedPolicy string
	if err == nil {
		schedPolicy = algorithm
	} else {
		schedPolicy = "DEFAULT_SCHEDPOLICY"
	}

	schd := Scheduler{
		ClusterClients:  cm.Cluster_kubeClients,
		ClusterInfoList: *clusterInfoList,
		//NodeInfoMap:     nodeInfoMap,
		ClusterInfoMap: clusterInfoMap,
		SchedPolicy:    schedPolicy,
	}

	return &schd
}

/*
func (sched *Scheduler) NewScoreTable() *scoretable.ScoreTable {
	var score_table scoretable.ScoreTable

	for _, cluster_info := range *sched.ClusterInfoList {
		var node_score_list scoretable.NodeScoreList

		// create node_score_list
		for _, node_info := range cluster_info.Nodes {
			node_score := &scoretable.NodeScore{
				Name:  node_info.NodeName,
				Score: 0,
			}
			node_score_list = append(node_score_list, *node_score)
		}
		cluster_score := &scoretable.ClusterScore{
			Cluster:       cluster_info.ClusterName,
			NodeScoreList: node_score_list,
			Score:         0,
		}
		score_table = append(score_table, *cluster_score)
	}

	return &score_table
}
*/

func (sched *Scheduler) Scheduling(deployment *v1alpha1.HCPDeployment) []v1alpha1.Target {

	fmt.Println("[scheduling start]")
	schedule_type := sched.SchedPolicy
	schedule_type = "NodeResourcesBalancedAllocation"

	fmt.Println("=> algorithm :", schedule_type)
	var num int32 = 2
	replicas := num

	var scheduling_result []v1alpha1.Target
	var cnt int32 = 0

	sched.SchedulingResource = test.TestPodsNodeResourcesBalancedAllocation[0]
	// set schedulingResource
	//sched.SchedulingResource = newPodFromHCPDeployment(deployment)

	//sched.Filtering("CheckNodeUnschedulable")
	for i := 0; i < int(replicas); i++ {
		sched.Scoring(schedule_type)
		best_cluster := sched.getMaxScoreCluster()
		fmt.Println(best_cluster)
		if best_cluster != "" {
			// if target !=  {
			if sched.updateSchedulingResult(best_cluster) {
				cnt += 1
				fmt.Printf("[Scheduling] %d/%d pod / TargetCluster : %s\n", i+1, replicas, best_cluster)
				fmt.Println()
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
		sched.printSchedulingResult()
		return scheduling_result
	} else {
		fmt.Println("Scheduling failed")
		return nil
	}

}

func (sched *Scheduler) getMaxScoreCluster() string {
	var max_score int32 = 0
	var best_cluster string = ""
	_ = best_cluster

	for key, value := range sched.ClusterInfoMap {
		//fmt.Println((*sched.ClusterInfoMap[key]).ClusterScore)
		if (*sched.ClusterInfoMap[key]).ClusterScore >= int32(max_score) {
			max_score = (*sched.ClusterInfoMap[key]).ClusterScore
			best_cluster = value.ClusterName
		}
	}

	return "test_cluster2"
}

func (sched *Scheduler) printSchedulingResult() {
	fmt.Println("========scheduling result========")
	targets := sched.SchedulingResult
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

func (sched *Scheduler) updateSchedulingResult(cluster string) bool {

	for i, target := range sched.SchedulingResult {
		// 이미 target cluster 목록에 cluster가 있는 경우
		if target.Cluster == cluster {
			// replicas 개수 증가
			temp := *target.Replicas
			temp += 1
			target.Replicas = &temp
			sched.SchedulingResult[i] = target
			return true
		}
	}

	// target cluster 목록에 cluster가 없는 경우

	// replicas 개수 1로 설정
	var new_target v1alpha1.Target
	new_target.Cluster = cluster
	var one int32 = 1
	new_target.Replicas = &one
	sched.SchedulingResult = append(sched.SchedulingResult, new_target)
	return true
}

func (sched *Scheduler) getTotalNumNodes() int {
	clusterinfoList := sched.ClusterInfoList
	cnt := 0

	for _, clusterinfo := range clusterinfoList {
		cnt += len(clusterinfo.Nodes)
	}

	return cnt
}

func (sched *Scheduler) Filtering(algorithm string) {
	var pod = &sched.SchedulingResource
	switch algorithm {
	case "CheckNodeUnschedulable":
		for i, _ := range sched.ClusterInfoList {
			fmt.Println(sched.ClusterInfoList[i].ClusterName)
			fmt.Println("==before filtering==")
			fmt.Println(sched.ClusterInfoList[i].Nodes)
			predicates.CheckNodeUnschedulable(*pod, &sched.ClusterInfoList[i])
			fmt.Println("==after filtering==")
			fmt.Println(&sched.ClusterInfoList[i].Nodes)
		}
	}
}

func (sched *Scheduler) Scoring(algorithm string) {

	fmt.Println("[scoring start]")
	var pod = &sched.SchedulingResource
	var score int32

	//clusterInfoMap := resourceinfo.CreateClusterInfoMap(&sched.ClusterInfoList)

	switch algorithm {

	case "Affinity":
		for _, clusterinfo := range sched.ClusterInfoList {
			fmt.Println("==>", clusterinfo.ClusterName)
			score = 0
			for _, node := range clusterinfo.Nodes {
				var node_score int32 = priorities.NodeAffinity(*pod, node.Node)
				if node_score == -1 {
					fmt.Println("fail to scoring node")
					return
				} else {
					node.NodeScore = node_score
					fmt.Println(node.NodeName, "score :", node_score)
					score += node_score
				}
			}
			sched.ClusterInfoMap[clusterinfo.ClusterName].ClusterScore = score
			fmt.Println("*", clusterinfo.ClusterName, "total score :", score)
		}
	case "TaintToleration":
		var node_score int32
		var result []int32

		// Get intolerable taints count
		for _, clusterinfo := range sched.ClusterInfoList {
			for _, node := range clusterinfo.Nodes {
				node_score = priorities.TaintToleration(*pod, node.Node)
				if node_score == -1 {
					fmt.Println("fail to scoring node")
					return
				} else {
					node.NodeScore = node_score
					result = append(result, node_score)
				}
			}
		}

		// sort intolerable taints count and get max value
		sort.Slice(result, func(i, j int) bool {
			return result[i] > result[j]
		})
		max := result[0]

		// scoring - normalize for intolerable taints count
		fmt.Println("[REAL SCORE]")
		for _, clusterinfo := range sched.ClusterInfoList {
			fmt.Println("==>", clusterinfo.ClusterName)
			score = 0
			for _, node := range clusterinfo.Nodes {
				if node.NodeScore == 0 {
					node.NodeScore = int32(scoretable.MaxNodeScore)
				} else {
					node.NodeScore = int32(100 * ((float32(max) - float32(node.NodeScore)) / float32(max)))
				}
				score += node.NodeScore
				fmt.Println("===>", node.NodeName, node.NodeScore)
			}

			if int32(len(clusterinfo.Nodes)) > 0 {
				// clusterInfoMap[clusterinfo.ClusterName].ClusterScore = score / int32(len((*clusterinfo).Nodes))
				sched.ClusterInfoMap[clusterinfo.ClusterName].ClusterScore = score
				fmt.Println("*", clusterinfo.ClusterName, "total score :", sched.ClusterInfoMap[clusterinfo.ClusterName].ClusterScore)
			} else {
				sched.ClusterInfoMap[clusterinfo.ClusterName].ClusterScore = 0
				fmt.Println("*", clusterinfo.ClusterName, "total score :", sched.ClusterInfoMap[clusterinfo.ClusterName].ClusterScore)
			}
		}
	case "NodeResourcesBalancedAllocation":
		// 현재 pod이 배치 된 후, CPU와 Memory 사용률이 균형을 검사
		for _, clusterinfo := range sched.ClusterInfoList {
			fmt.Println("==>", clusterinfo.ClusterName)
			score = 0
			for _, node := range clusterinfo.Nodes {
				var node_score int32 = int32(priorities.NodeResourcesBalancedAllocation(*pod, node.Node))
				if node_score == -1 {
					fmt.Println("fail to scoring node")
					return
				} else {
					node.NodeScore = node_score
					fmt.Println(node.NodeName, "score :", node_score)
					score += node_score
				}
			}

			//fmt.Println(sched.ClusterInfoMap)
			(*sched.ClusterInfoMap[clusterinfo.ClusterName]).ClusterScore = score
			fmt.Println("*", clusterinfo.ClusterName, "total score :", score)
			fmt.Println()
		}
	case "ImageLocality":
		for _, clusterinfo := range sched.ClusterInfoList {
			fmt.Println("==>", clusterinfo.ClusterName)
			score = 0
			for _, node := range clusterinfo.Nodes {
				fmt.Println(node.ImageStates)
				var node_score int32 = int32(priorities.ImageLocality(*pod, node, sched.getTotalNumNodes()))
				if node_score == -1 {
					fmt.Println("fail to scoring node")
					return
				} else {
					node.NodeScore = node_score
					fmt.Println(node.NodeName, "score :", node_score)
					score += node_score
				}
			}
			sched.ClusterInfoMap[clusterinfo.ClusterName].ClusterScore = score
			fmt.Println("*", clusterinfo.ClusterName, "total score :", score)
		}
	}

}

func (sched *Scheduler) scheduleOne(ctx context.Context) {

}
