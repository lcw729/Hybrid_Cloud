package main

import (
	monitoringEngine "Hybrid_Cloud/hcp-monitoring-engine/metricCollector"
	"fmt"
)

func main() {
	fmt.Println("MetricCollector Called")
	monitoringEngine.MetricCollector()
}
