package metricCollector

import (
	"fmt"
)

type MetricPoint struct {
	Timestamp             string `json:"timestamp"`
	CPUUsageNanoCores     string `json:"cpu_usage_bytes"`
	MemoryUsageBytes      string `json:"memory_usage_bytes"`
	MemoryAvailableBytes  string `json:"memeory_available_bytes"`
	MemoryWorkingSetBytes string `json:"memory_workingset_bytes"`
	NetworkRxBytes        string `json:"network_rx_bytes"`
	NetworkTxBytes        string `json:"network_tx_bytes"`
	FsAvailableBytes      string `json:"fs_available_bytes"`
	FsCapacityBytes       string `json:"fs_capacity_bytes"`
	FsUsedBytes           string `json:"fs_usage_bytes"`
	NetworkLatency        string `json:"network_latency"`
}

type Monitoring struct {
	Name       string      `json:"name"`
	Namespace  string      `json:"name_space"`
	MP         MetricPoint `json:"mp"`
	Containers []int       `json:"containers"`
}

func MetricCollector() {
	// mp := &MetricPoint{
	// 	Timestamp:             "2021-08-07T12:03:12Z",
	// 	CPUUsageNanoCores:     "2",
	// 	MemoryUsageBytes:      "283MiB",
	// 	MemoryAvailableBytes:  "3412GiB",
	// 	MemoryWorkingSetBytes: "122MiB",
	// 	NetworkRxBytes:        "23MiB",
	// 	NetworkTxBytes:        "312MiB",
	// 	FsAvailableBytes:      "412MiB",
	// 	FsCapacityBytes:       "2MiB",
	// 	FsUsedBytes:           "12MiB",
	// 	NetworkLatency:        "0.12mseconds",
	// }
	// pmp := &Monitoring{
	// 	Name:       "metric-test",
	// 	Namespace:  "hybrid",
	// 	MP:         *mp,
	// 	Containers: nil,
	// }

	// podMetricPoints, err := json.MarshalIndent(pmp, "", " ")
	// if err != nil {
	// 	panic(err)
	// }
	fmt.Println("MetricCollector Called")
	fmt.Println(" ")
	// fmt.Println(string(podMetricPoints))
	// fmt.Println(" ")
	fmt.Println("---------------------------------------")
	fmt.Println("Time Stamp: 2021-08-07T12:03:12Z")
	fmt.Println("CPU Usage Nano Cores : 2")
	fmt.Println("Memmory Usage Bytes : 283MiB")
	fmt.Println("Memmory Available Bytes : 3412GiB")
	fmt.Println("Memmory Working Bytes : 122MiB")
	fmt.Println("Network Rx Bytes : 23MiB")
	fmt.Println("Network Tx Bytes : 312MiB")
	fmt.Println("Fs Available Bytes : 412MiB")
	fmt.Println("Fs Capacity Bytes : 2MiB")
	fmt.Println("Fs Used Bytes : 12MiB")
	fmt.Println("Network Latency : 12MS")
	fmt.Println("Pod Name : metric-test")
	fmt.Println("Name Space : hybrid")
	fmt.Println("Container : keti-container")
	fmt.Println("---------------------------------------")
	fmt.Println("")
	fmt.Println("Send MultiMetric")
}
