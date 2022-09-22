package algorithm

/*
import (
	"Hybrid_Cloud/hcp-analytic-engine/pkg/handler"
	"Hybrid_Cloud/hcp-analytic-engine/pkg/metric"
	policy "Hybrid_Cloud/hcp-resource/hcppolicy"
	"Hybrid_Cloud/util/clusterManager"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"k8s.io/klog"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

var AlgorithmMap = map[string]func() bool{
	"DRF":      DRF,
	"Affinity": Affinity,
}

var TargetCluster = make(map[string]*fedv1b1.KubeFedCluster)
var cm, _ = clusterManager.NewClusterManager()

func Calculate_WatchingLevel(index int, podNum int, clusterName string) (bool, string, string, error) {

	klog.Info("----------------------------------------------------------")
	klog.Info("[step 1] Get Policy - watching level & warning level")
	klog.Info("< Watching Level >")
	watching_level := policy.GetWatchingLevel(*cm.HCPPolicy_Client)
	for _, level := range watching_level.Levels {
		fmt.Printf("watching level %s : %s\n", level.Type, level.Value)
	}
	//	각 클러스터의 watching level 계산하고 warning level 초과 시 targetCluster에 추가
	warning_level, err := strconv.Atoi(policy.GetWarningLevel(*cm.HCPPolicy_Client).Value)
	if err != nil {
		return false, "", "", err
	}

	klog.Info("< Warning  Level >")
	fmt.Printf("warning level: %d\n", warning_level)
	klog.Info("----------------------------------------------------------")
	klog.Info("[step 2] Get MultiMetric")
	var watchingLevel = 0
	var jsonarray metric.PodMetric
	jsonByteArray := handler.GetResource(podNum, clusterName, "_")
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		//panic(err)
		return false, "", "", err
	}
	klog.Info("Print Cluster : ", clusterName)
	klog.Info("")

	//파드 개수를 가져와서 for문의 변수로 넣어주어야 함

	klog.Info("///////////////////////////////////////////////////")
	klog.Info("Pod", index+1, " Metric Information")
	klog.Info("PodMetric         : ", jsonarray.Podmetrics[index].Pod)
	klog.Info("Node", jsonarray.Podmetrics[index].Node)
	klog.Info("Time              : ", jsonarray.Podmetrics[index].Time)
	klog.Info("Cpu               : ", jsonarray.Podmetrics[index].CPU.CPUUsageNanoCores)
	klog.Info("Cluster           : ", jsonarray.Podmetrics[index].Cluster)
	klog.Info("Memory            : ", jsonarray.Podmetrics[index].Memory)
	klog.Info("Namespace         : ", jsonarray.Podmetrics[index].Namespace)
	klog.Info("///////////////////////////////////////////////////")
	klog.Info("----------------------------------------------------------")
	klog.Info("[step 3] Calculate watching level")
	CpuUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[index].CPU.CPUUsageNanoCores, "n"), 32)
	if err != nil {
		// handle error
		klog.Error(err)
		os.Exit(2)
	}
	CpuUsage = CpuUsage / 1000000000
	if CpuUsage != 0 {
		klog.Info("CpuUsage          : ", CpuUsage)
		switch {
		case CpuUsage <= 0.2:
			watchingLevel = 1
			klog.Info("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage <= 0.4:
			watchingLevel = 2
			klog.Info("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage <= 0.6:
			watchingLevel = 3
			klog.Info("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage <= 0.8:
			watchingLevel = 4
			klog.Info("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage <= 1:
			watchingLevel = 5
			klog.Info("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
		case CpuUsage > 1:
			watchingLevel = 5
			klog.Info("Pod Name ", jsonarray.Podmetrics[index].Pod, "watching Level is ", watchingLevel)
			klog.Info("Watching Level is over 5")
			klog.Info("!!!!!!!!!!!!!!!", jsonarray.Podmetrics[index].Pod, "!!!!!!!!!!!!!!!!!!!")
		}
	} else {
		klog.Info("This Pod use 0 Cpu nanocores")
		return false, "", "", nil
	}

	if watchingLevel >= warning_level {
		klog.Info("watching level is over warning level!!!!!!!")
		klog.Info("-----------------------------------------")
		return true, jsonarray.Podmetrics[index].Pod, jsonarray.Podmetrics[index].Namespace, nil // 초과  -- HPA, VPA 수행
	} else {
		klog.Info("-----------------------------------------")
		return false, "", "", nil
	}

}

func PastWatchingLevelCalculator() (bool, error) {
	klog.Info("-----------------------------------------")
	klog.Info("[step 1] Get Policy - watching level & warning level")
	klog.Info("< Watching Level >")
	watching_level := policy.GetWatchingLevel(*cm.HCPPolicy_Client)
	for _, level := range watching_level.Levels {
		fmt.Printf("watching level %s : %s\n", level.Type, level.Value)
	}
	// 각 클러스터의 watching level 계산하고 warning level 초과 시 targetCluster에 추가
	warning_level, err := strconv.Atoi(policy.GetWarningLevel(*cm.HCPPolicy_Client).Value)
	if err != nil {
		return false, err
	}
	klog.Info("< Warning  Level >")
	fmt.Printf("warning level: %d\n", warning_level)
	klog.Info("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	klog.Info("[step 3] Calculate watching level")
	result := 3
	fmt.Printf(">>> Result <<< %d level\n", result)

	if result >= warning_level {
		klog.Info("watching level is over warning level!!!!!!!")
		klog.Info("-----------------------------------------")
		return true, nil // 초과  -- HPA, VPA 수행
	} else {
		klog.Info("-----------------------------------------")
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
// 	// klog.Info("[step 5] Start ResourceConfiguration")
// 	// for index, cluster := range targetCluster {
// 	// 	klog.Info("Index : ", index, "\nClusterId : ", cluster.ClusterId, "\nClusterName : ", cluster.ClusterName)
// 	// }
// }

// 최적 배치 알고리즘
func Affinity() bool {
	klog.Info("---------------------------------------------------------------")
	klog.Info("Affinity Calculator Called")
	klog.Info("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	klog.Info("[step 3-1] Start analysis Resource Affinity")
	score_table["hcp-cluster"] = 50.0
	score_table["a-cluster"] = 40.0
	score_table["b-cluster"] = 20.0
	score_table["c-cluster"] = 10.0
	klog.Info("[step 3-2] Send analysis result to Scheduler [Target Cluster]")
	klog.Info("---------------------------------------------------------------")
	return true
}

func DRF() bool {
	klog.Info("DRF Math operation Called")
	klog.Info("-----------------------------------------")
	klog.Info("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	klog.Info("-----------------------------------------")
	klog.Info("[step 3-1] Handling Math Operation")
	klog.Info("[step 3-2] Search Pod Fit Resources")
	klog.Info("[step 3-3] Schedule Decision")
	klog.Info("---------------------------------------------------------------")
	return true
}

/*
// 배치 알고리즘에 설정 값에 따라 알고리즘 변경
func OptimalArrangementAlgorithm() map[string]float32 {
	klog.Info("[step 1] Get Policy - algorithm")
	algo := policy.GetAlgorithm()
	if algo != "" {
		klog.Info(algo)
		switch algo {
		case "DRF":
			return DRF()
		case "Affinity":
			return Affinity()
		default:
			klog.Info("there is no such algorithm.")
			return false
		}
	} else {
		klog.Info("there is no such algorithm.")
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
