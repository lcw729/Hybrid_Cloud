package analyticEngine

import (
	"Hybrid_Cluster/hcp-analytic-engine/util"
	monitoringEngine "Hybrid_Cluster/hcp-monitoring-engine/metricCollector"
	"Hybrid_Cluster/hcp-scheduler/pkg/policy"
	"fmt"
)

func WatchingLevelCalculator() []*util.ClusterInfo {
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 3] Get Policy - watching level & warning level")
	fmt.Println("< Watching Level > \n", policy.GetWatchingLevel())
	// 각 클러스터의 watching level 계산하고 warning level 초과 시 targetCluster에 추가
	fmt.Println("< Warning  Level > \n", policy.GetWarningLevel().Value)
	fmt.Println("[step 4] Calculate watching level")
	var targetCluster []*util.ClusterInfo

	targetCluster = append(targetCluster, &util.ClusterInfo{
		ClusterId:   1,
		ClusterName: "cluster1",
	})
	targetCluster = append(targetCluster, &util.ClusterInfo{
		ClusterId:   2,
		ClusterName: "cluster2",
	})
	targetCluster = append(targetCluster, &util.ClusterInfo{
		ClusterId:   3,
		ClusterName: "cluster3",
	})
	return targetCluster
}

func ResourceConfigurationManagement() {
	targetCluster := WatchingLevelCalculator()
	fmt.Println("[step 5] Start ResourceConfiguration")
	for index, cluster := range targetCluster {
		fmt.Println("Index : ", index, "\nClusterId : ", cluster.ClusterId, "\nClusterName : ", cluster.ClusterName)
	}
}

func HybridctlAnalyticEngine() {

	fmt.Println("Hybridctl Analytic Engine Called")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 1] Get MultiMetric")
	monitoringEngine.MetricCollector()
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 2] Start analysis to find target cluster")
	fmt.Println("[step 3] Send result to API Server")
	fmt.Println("---------------------------------------------------------------")

}

func ResourceConfigurationManagement() {
	fmt.Println("Resource Extenstion Analytic Engine Called")
	fmt.Println("[step 4] Start analysis to find target cluster")
	fmt.Println("[step 5] Call HPA or VPA")

}

func AffinityCalculator() {
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("Affinity Calculator Called")
	fmt.Println("[step 2-1] Get MultiMetric")
	monitoringEngine.MetricCollector()
	fmt.Println("[step 2-2] Start analysis Resource Affinity")
	fmt.Println("[step 2-3] Send analysis result to Scheduler [Target Cluster]")
	fmt.Println("---------------------------------------------------------------")
}

// func PodNodeSelector() {
// 	fmt.Println("---------------------------------------------------------------")
// 	fmt.Println("PodNodeSelector Called")
// 	fmt.Println("[step 2-1] Get MultiMetric")
// 	monitoringEngine.MetricCollector()
// 	fmt.Println("[step 2-2] Start analysis Resource Affinity")
// 	fmt.Println("[step 2-3] Send analysis result to Scheduler [Target Cluster]")
// 	fmt.Println("---------------------------------------------------------------")
// }

func DRF() {
	fmt.Println("DRF Math operation Called")
	fmt.Println("[step 3-1] Handling Math Operation")
	fmt.Println("[step 3-2] Search Pod Fit Resources")
	fmt.Println("[step 3-3] Schedule Decision")
	fmt.Println("---------------------------------------------------------------")
}
