package algorithm

import (
	policy "Hybrid_Cluster/hcp-resource/hcppolicy"
	clusterManager "Hybrid_Cluster/util/clusterManager"
	"fmt"

	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

var AlgorithmMap = map[string]func() bool{
	"DRF":      DRF,
	"Affinity": Affinity,
}

var TargetCluster = make(map[string]*fedv1b1.KubeFedCluster)
var watching_level = policy.GetWatchingLevel()
var warning_level = policy.GetWarningLevel()

func WatchingLevelCalculator() {
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 2] Get Policy - watching level & warning level")
	fmt.Println("< Watching Level > \n", watching_level)
	// 각 클러스터의 watching level 계산하고 warning level 초과 시 targetCluster에 추가
	fmt.Println("< Warning  Level > \n", warning_level)
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 3] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	fmt.Println("[step 4] Calculate watching level")

	cm := clusterManager.NewClusterManager()
	for _, cluster := range cm.Cluster_list.Items {
		fmt.Println(cluster.Name)
		TargetCluster[cluster.Name] = &cluster
	}
	fmt.Println(TargetCluster)
	// cluster := &util.ClusterInfo{
	// 	ClusterId:   1,
	// 	ClusterName: "cluster1",
	// }
	// if !appendTargetCluster(cluster) {
	// 	fmt.Printf("%d exist already\n", cluster.ClusterId)
	// } else {
	// 	fmt.Println("ok")
	// }

	// cluster = &util.ClusterInfo{
	// 	ClusterId:   2,
	// 	ClusterName: "cluster2",
	// }
	// if !appendTargetCluster(cluster) {
	// 	fmt.Printf("%d exist already\n", cluster.ClusterId)
	// } else {
	// 	fmt.Println("ok")
	// }
}

// func appendTargetCluster(cluster *util.ClusterInfo) bool {
// 	var check bool = false
// 	for _, c := range util.TargetCluster {
// 		if c.ClusterId == cluster.ClusterId {
// 			check = true
// 			break
// 		}
// 	}
// 	if !check {
// 		util.TargetCluster = append(util.TargetCluster, cluster)
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func ResourceConfigurationManagement() {
// 	// targetCluster := WatchingLevelCalculator()
// 	WatchingLevelCalculator()
// 	// fmt.Println("[step 5] Start ResourceConfiguration")
// 	// for index, cluster := range targetCluster {
// 	// 	fmt.Println("Index : ", index, "\nClusterId : ", cluster.ClusterId, "\nClusterName : ", cluster.ClusterName)
// 	// }
// }

// 최적 배치 알고리즘
func Affinity() bool {
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Affinity Calculator Called")
	fmt.Println("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	fmt.Println("[step 3-1] Start analysis Resource Affinity")
	score_table["a-cluster"] = 30.0
	score_table["b-cluster"] = 40.0
	score_table["c-cluster"] = 20.0
	score_table["b"] = 50.0
	fmt.Println("[step 3-2] Send analysis result to Scheduler [Target Cluster]")
	fmt.Println("---------------------------------------------------------------")
	return true
}

func DRF() bool {
	fmt.Println("DRF Math operation Called")
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 3-1] Handling Math Operation")
	fmt.Println("[step 3-2] Search Pod Fit Resources")
	fmt.Println("[step 3-3] Schedule Decision")
	fmt.Println("---------------------------------------------------------------")
	return true
}

/*
// 배치 알고리즘에 설정 값에 따라 알고리즘 변경
func OptimalArrangementAlgorithm() map[string]float32 {
	fmt.Println("[step 1] Get Policy - algorithm")
	algo := policy.GetAlgorithm()
	if algo != "" {
		fmt.Println(algo)
		switch algo {
		case "DRF":
			return DRF()
		case "Affinity":
			return Affinity()
		default:
			fmt.Println("there is no such algorithm.")
			return false
		}
	} else {
		fmt.Println("there is no such algorithm.")
		return false
	}
}
*/

/*
// 가장 점수가 높은 Cluster, Node 확인
func OptimalNodeSelector() (*util.Cluster, *util.NodeScore) {
	max := 0
	cluster := util.ScoreTable[0]
	node := util.ScoreTable[0].Nodes[0]
	for _, c := range util.ScoreTable {
		for _, n := range c.Nodes {
			if int(n.Score) > max {
				max = int(n.Score)
				cluster = c
				node = n
			}
		}
	}
	return cluster, node
}
*/
