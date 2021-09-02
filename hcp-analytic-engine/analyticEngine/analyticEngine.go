package analyticEngine

import (
	monitoringEngine "Hybrid_Cluster/hcp-monitoring-engine/metricCollector"
	"encoding/json"
	"fmt"
)

type ResourceWeightResult struct {
	Cluster1 int `json:"Cluster1"`
	Cluster2 int `json:"Cluster2"`
	Cluster3 int `json:"Cluster3"`
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

func ResourceExtension() {
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 2] Get MultiMetric")
	monitoringEngine.MetricCollector()
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 3] Calculate resource weight")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 4] Send weight calculation result to Scheduler (Resource Balancing Controller)")
	fmt.Println("--Resource Weight Result--")
	result := ResourceWeightResult{
		Cluster1: 30,
		Cluster2: 60,
		Cluster3: 10,
	}
	r, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(r))
	fmt.Println("---------------------------------------------------------------")

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
