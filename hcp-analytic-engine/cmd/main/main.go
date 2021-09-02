package main

import (
	analyticEngine "Hybrid_Cluster/hcp-analytic-engine/analyticEngine"
	"fmt"
)

func main() {
	fmt.Println("---------------------------------")
	//analyticEngine.HybridctlAnalyticEngine()
	// analyticEngine.AffinityCalculator()
	// analyticEngine.DRF()
	analyticEngine.ResourceConfigurationManagement()
	// analyticEngine.ResourceExtension()
	fmt.Println("---------------------------------")
	// http.HandleFunc("/hybridctlAnalyticEngine", HybridctlAnalyticEngine)
	// http.HandleFunc("/ResourceExtensione", ResourceExtension)
	// http.HandleFunc("/ResourceConfigurationManagement", ResourceConfigurationManagement)
	// http.HandleFunc("/AffinityCalculator", AffinityCalculator)
	// http.HandleFunc("/DRF", DRF)
	// http.ListenAndServe(":8090", nil)

}
