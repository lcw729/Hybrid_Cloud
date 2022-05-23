package metric

import "time"

type PodMetric struct {
	Podmetrics []struct {
		Time      time.Time `json:"time"`
		Cluster   string    `json:"cluster"`
		Namespace string    `json:"namespace"`
		Node      string    `json:"node"`
		Pod       string    `json:"pod"`
		CPU       struct {
			CPUUsageNanoCores string `json:"CPUUsageNanoCores"`
		} `json:"cpu"`
		Memory struct {
			MemoryAvailableBytes  string `json:"MemoryAvailableBytes"`
			MemoryUsageBytes      string `json:"MemoryUsageBytes"`
			MemoryWorkingSetBytes string `json:"MemoryWorkingSetBytes"`
		} `json:"memory"`
		Fs struct {
			FsAvailableBytes string `json:"FsAvailableBytes"`
			FsCapacityBytes  string `json:"FsCapacityBytes"`
			FsUsedBytes      string `json:"FsUsedBytes"`
		} `json:"fs"`
		Network struct {
			NetworkRxBytes string `json:"NetworkRxBytes"`
			NetworkTxBytes string `json:"NetworkTxBytes"`
		} `json:"network"`
	} `json:"podmetrics"`
}
