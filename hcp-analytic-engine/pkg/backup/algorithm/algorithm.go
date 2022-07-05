package algorithm

import (
	"Hybrid_Cloud/hcp-analytic-engine/pkg/handler"
	"Hybrid_Cloud/hcp-analytic-engine/pkg/metric"
	policy "Hybrid_Cloud/hcp-resource/hcppolicy"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

var AlgorithmMap = map[string]func() bool{
	"DRF":      DRF,
	"Affinity": Affinity,
}

var TargetCluster = make(map[string]*fedv1b1.KubeFedCluster)

func Calculate_WatchingLevel(index int, podNum int, clusterName string) (bool, string, string, error) {

	fmt.Println("----------------------------------------------------------")
	fmt.Println("[step 1] Get Policy - watching level & warning level")
	fmt.Println("< Watching Level >")
	watching_level := policy.GetWatchingLevel()
	for _, level := range watching_level.Levels {
		fmt.Printf("watching level %s : %s\n", level.Type, level.Value)
	}
	//	각 클러스터의 watching level 계산하고 warning level 초과 시 targetCluster에 추가
	warning_level, err := strconv.Atoi(policy.GetWarningLevel().Value)
	if err != nil {
		return false, "", "", err
	}

	fmt.Println("< Warning  Level >")
	fmt.Printf("warning level: %d\n", warning_level)
	fmt.Println("----------------------------------------------------------")
	fmt.Println("[step 2] Get MultiMetric")
	var watchingLevel = 0
	var jsonarray metric.PodMetric
	jsonByteArray := handler.GetResource(podNum, clusterName, "_")
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		//panic(err)
		return false, "", "", err
	}
	fmt.Println("Print Cluster : ", clusterName)
	fmt.Println("")

	//파드 개수를 가져와서 for문의 변수로 넣어주어야 함

	fmt.Println("///////////////////////////////////////////////////")
	fmt.Println("Pod", index+1, " Metric Information")
	fmt.Println("PodMetric         : ", jsonarray.Podmetrics[index].Pod)
	fmt.Println("Node", jsonarray.Podmetrics[index].Node)
	fmt.Println("Time              : ", jsonarray.Podmetrics[index].Time)
	fmt.Println("Cpu               : ", jsonarray.Podmetrics[index].CPU.CPUUsageNanoCores)
	fmt.Println("Cluster           : ", jsonarray.Podmetrics[index].Cluster)
	fmt.Println("Memory            : ", jsonarray.Podmetrics[index].Memory)
	fmt.Println("Namespace         : ", jsonarray.Podmetrics[index].Namespace)
	fmt.Println("///////////////////////////////////////////////////")
	fmt.Println("----------------------------------------------------------")
	fmt.Println("[step 3] Calculate watching level")
	CpuUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[index].CPU.CPUUsageNanoCores, "n"), 32)
	if err != nil {
		// handle error
		fmt.Println(err)
		os.Exit(2)
	}
	CpuUsage = CpuUsage / 1000000000
	if CpuUsage != 0 {
		fmt.Println("CpuUsage          : ", CpuUsage)
		switch {
		case CpuUsage <= 0.2:
			watchingLevel = 1
			fmt.Println("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage <= 0.4:
			watchingLevel = 2
			fmt.Println("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage <= 0.6:
			watchingLevel = 3
			fmt.Println("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage <= 0.8:
			watchingLevel = 4
			fmt.Println("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage <= 1:
			watchingLevel = 5
			fmt.Println("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage > 1:
			watchingLevel = 5
			fmt.Println("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
			fmt.Println("Watching Level is over 5")
			fmt.Println("!!!!!!!!!!!!!!!", jsonarray.Podmetrics[index].Pod, "!!!!!!!!!!!!!!!!!!!")
		}
	} else {
		fmt.Println("This Pod use 0 Cpu nanocores")
		return false, "", "", nil
	}

	if watchingLevel >= warning_level {
		fmt.Println("watching level is over warning level!!!!!!!")
		fmt.Println("-----------------------------------------")
		return true, jsonarray.Podmetrics[index].Pod, jsonarray.Podmetrics[index].Namespace, nil // 초과  -- HPA, VPA 수행
	} else {
		fmt.Println("-----------------------------------------")
		return false, "", "", nil
	}

}

func PastWatchingLevelCalculator() (bool, error) {
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 1] Get Policy - watching level & warning level")
	fmt.Println("< Watching Level >")
	watching_level := policy.GetWatchingLevel()
	for _, level := range watching_level.Levels {
		fmt.Printf("watching level %s : %s\n", level.Type, level.Value)
	}
	// 각 클러스터의 watching level 계산하고 warning level 초과 시 targetCluster에 추가
	warning_level, err := strconv.Atoi(policy.GetWarningLevel().Value)
	if err != nil {
		return false, err
	}
	fmt.Println("< Warning  Level >")
	fmt.Printf("warning level: %d\n", warning_level)
	fmt.Println("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	fmt.Println("[step 3] Calculate watching level")
	result := 3
	fmt.Printf(">>> Result <<< %d level\n", result)

	if result >= warning_level {
		fmt.Println("watching level is over warning level!!!!!!!")
		fmt.Println("-----------------------------------------")
		return true, nil // 초과  -- HPA, VPA 수행
	} else {
		fmt.Println("-----------------------------------------")
		return false, nil
	}

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
	score_table["hcp-cluster"] = 50.0
	score_table["a-cluster"] = 40.0
	score_table["b-cluster"] = 20.0
	score_table["c-cluster"] = 10.0
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
