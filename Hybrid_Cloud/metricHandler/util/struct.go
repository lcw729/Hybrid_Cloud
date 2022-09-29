package util

type PodMetricList struct {
	Items []PodMetric `json:"podmetrics"`
}
type NodeMetricList struct {
	Items []NodeMetric `json:"nodemetrics"`
}
type PodMetric struct {
	Time      string        `json:"time"`
	Cluster   string        `json:"cluster"`
	Namespace string        `json:"namespace"`
	Node      string        `json:"node"`
	Pod       string        `json:"pod"`
	Cpu       CpuMetric     `json:"cpu"`
	Memory    MemoryMetric  `json:"memory"`
	Fs        FsMetric      `json:"fs"`
	Network   NetworkMetric `json:"network"`
}
type NodeMetric struct {
	Time    string        `json:"time"`
	Cluster string        `json:"cluster"`
	Node    string        `json:"node"`
	Cpu     CpuMetric     `json:"cpu"`
	Memory  MemoryMetric  `json:"memory"`
	Fs      FsMetric      `json:"fs"`
	Network NetworkMetric `json:"network"`
}

type CpuMetric struct {
	CPUUsageNanoCores string `json:"CPUUsageNanoCores"`
}
type MemoryMetric struct {
	MemoryAvailableBytes  string `json:"MemoryAvailableBytes"`
	MemoryUsageBytes      string `json:"MemoryUsageBytes"`
	MemoryWorkingSetBytes string `json:"MemoryWorkingSetBytes"`
}
type FsMetric struct {
	FsAvailableBytes string `json:"FsAvailableBytes"`
	FsCapacityBytes  string `json:"FsCapacityBytes"`
	FsUsedBytes      string `json:"FsUsedBytes"`
}
type NetworkMetric struct {
	NetworkRxBytes string `json:"NetworkRxBytes"`
	NetworkTxBytes string `json:"NetworkTxBytes"`
}
