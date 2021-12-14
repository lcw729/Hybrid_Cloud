// package main

// import (
// 	"Hybrid_Cluster/test/pkg/handler"
// 	"encoding/json"
// 	"fmt"
// 	"time"
// )

// func main() {
// 	fmt.Println("test module")
// 	var jsonarray PodMetric
// 	ns := "hcp"
// 	var podNum = []int{12, 9, 13}
// 	var clusterName = []string{"gke-cluster", "eks-cluster", "aks-cluster"}
// 	for i := 0; i < 3; i++ {
// 		jsonByteArray := handler.GetResource(podNum[i], ns, clusterName[i])
// 		stringArray := string(jsonByteArray[:])
// 		if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
// 			panic(err)
// 		}
// 		for j := 0; j < podNum[i]; j++ {
// 			fmt.Println("Pod", i+1, " Metric Information")
// 			fmt.Println("PodMetric: ", jsonarray.Podmetrics[i].Pod)
// 			fmt.Println("time: ", jsonarray.Podmetrics[i].Time)
// 			fmt.Println("cpu: ", jsonarray.Podmetrics[i].CPU)
// 			fmt.Println("cluster: ", jsonarray.Podmetrics[i].Cluster)
// 			fmt.Println("memory: ", jsonarray.Podmetrics[i].Memory)
// 			fmt.Println("namespace: ", jsonarray.Podmetrics[i].Namespace)
// 		}
// 	}

// 	// stringArray := string(jsonByteArray[:])
// 	// if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
// 	// 	panic(err)
// 	// }
// }

// type PodMetric struct {
// 	Podmetrics []struct {
// 		Time      time.Time `json:"time"`
// 		Cluster   string    `json:"cluster"`
// 		Namespace string    `json:"namespace"`
// 		Node      string    `json:"node"`
// 		Pod       string    `json:"pod"`
// 		CPU       struct {
// 			CPUUsageNanoCores string `json:"CPUUsageNanoCores"`
// 		} `json:"cpu"`
// 		Memory struct {
// 			MemoryAvailableBytes  string `json:"MemoryAvailableBytes"`
// 			MemoryUsageBytes      string `json:"MemoryUsageBytes"`
// 			MemoryWorkingSetBytes string `json:"MemoryWorkingSetBytes"`
// 		} `json:"memory"`
// 		Fs struct {
// 			FsAvailableBytes string `json:"FsAvailableBytes"`
// 			FsCapacityBytes  string `json:"FsCapacityBytes"`
// 			FsUsedBytes      string `json:"FsUsedBytes"`
// 		} `json:"fs"`
// 		Network struct {
// 			NetworkRxBytes string `json:"NetworkRxBytes"`
// 			NetworkTxBytes string `json:"NetworkTxBytes"`
// 		} `json:"network"`
// 	} `json:"podmetrics"`
// }

package main

import (
	"Hybrid_Cluster/test/pkg/handler"
	"encoding/json"
	"fmt"
	"time"
)

func main() {

	var cluster_list = []string{"gke-cluster", "eks-cluster", "aks-cluster"}

	var podNum = []int{13, 9, 13}

	for {
		for i := 0; i < len(cluster_list); i++ {
			Calculate_WatchingLevel(podNum[i], cluster_list[i])
		}
		time.Sleep(2 * time.Second)
	}
}

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

func Calculate_WatchingLevel(podNum int, clusterName string) {
	var jsonarray PodMetric
	ns := "hcp"
	jsonByteArray := handler.GetResource(podNum, ns, clusterName)
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		panic(err)
	}
	fmt.Println("----------------------------------------------------------")
	fmt.Println("Print Cluster : ", clusterName)
	fmt.Println("")

	for i := 0; i < podNum; i++ {
		fmt.Println("Pod", i+1, " Metric Information")
		fmt.Println("PodMetric: ", jsonarray.Podmetrics[i].Pod)
		fmt.Println("time: ", jsonarray.Podmetrics[i].Time)
		fmt.Println("cpu: ", jsonarray.Podmetrics[i].CPU)
		fmt.Println("cluster: ", jsonarray.Podmetrics[i].Cluster)
		fmt.Println("memory: ", jsonarray.Podmetrics[i].Memory)
		fmt.Println("namespace: ", jsonarray.Podmetrics[i].Namespace)
	}
	fmt.Println("------------------------------------------")
}
