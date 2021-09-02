package main

import (
	"Hybrid_Cluster/hcp-analytic-engine/analyticEngine"
	monitoringEngine "Hybrid_Cluster/hcp-monitoring-engine/metricCollector"
	"fmt"
)

func main() {
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 1] Get MultiMetric")
	monitoringEngine.MetricCollector()
	fmt.Println("-----------------------------------------")
	fmt.Println("[step 2] Get Policy - watching level")
	fmt.Println("----> Policy 1: setting the standard of watching level")
	fmt.Println("----> Policy 2: settting target level")
	fmt.Println("[step 3].Calculate watching level")
	fmt.Println("----> LEVEL : 4")
	fmt.Println("----> limit exceeded!!!!!")
	// fmt.Printf("----> When the watching level exceeds target level,\n the multi-metric information of the corresponding pod is transmitted to the analysis engine.\n")
	fmt.Println("-----------------------------------------")
	analyticEngine.ResourceConfigurationManagement()

}
