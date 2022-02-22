// package main

// import (
// 	"Hybrid_Cluster/test/pkg/handler"
// 	"encoding/json"
// 	"fmt"
// 	"time"
// )
package main

import (
	hcppolicyapis "Hybrid_Cluster/pkg/apis/hcppolicy/v1alpha1"
	hcppolicyv1alpha1 "Hybrid_Cluster/pkg/client/hcppolicy/v1alpha1/clientset/versioned"
	"Hybrid_Cluster/test/pkg/handler"
	"Hybrid_Cluster/util/clusterManager"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

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

func main() {
	cpu, _ := GetInitialSettingValue("max_cpu")
	mem, _ := GetInitialSettingValue("max_memory")
	extra, _ := GetInitialSettingValue("extra")

	fmt.Println(extra)
	cpu = cpu * (100 - extra) / 100
	mem = mem * (100 - extra) / 100
	fmt.Println(cpu, mem)
	// c := influx.InfluxDBClient("10.0.5.83", "31051", "root", "root")
	// c.Query("select \"CPUUsageNanoCores\", \"MemoryUsageBytes\", \"NetworkLatency\"  from (select * from Pods where \"cluster\"='"+cluster+"' and \"namespace\"='"+namespace+"' and \"pod\"=~/"+depname+"/ order by time DESC limit "+podnum+") order by time desc", "Metrics", "")
	/*
		var cluster_list = []string{"gke-cluster", "eks-cluster", "aks-cluster"}

		var podNum = []int{13, 9, 13}

		for {
			for i := 0; i < len(cluster_list); i++ {
				Calculate_WatchingLevel(podNum[i], cluster_list[i])
			}
			time.Sleep(2 * time.Second)
		}
	*/
}

func GetInitialSettingValue(typ string) (int, string) {
	policy := GetPolicy("initial-setting")
	policies := policy.Spec.Template.Spec.Policies
	for _, p := range policies {
		println(p.Type)
		if typ == "default_node_option" && p.Type == "default_node_option" {
			var value string
			value = p.Value
			if value == "" {
				fmt.Printf("ERROR: No %s Value\n", typ)
			} else {
				return -1, value
			}
		} else if p.Type == typ {
			var value int
			value, err := strconv.Atoi(p.Value)
			if err != nil {
				fmt.Printf("ERROR: No %s Value\n", typ)
			}
			return value, ""
		}
	}
	fmt.Printf("ERROR: No Such Type %s\n", typ)
	return -1, ""
}

func GetPolicy(policy_name string) *hcppolicyapis.HCPPolicy {
	cm := clusterManager.NewClusterManager()

	c, err := hcppolicyv1alpha1.NewForConfig(cm.Host_config)
	if err != nil {
		klog.Info(err)
	}
	policy, err := c.HcpV1alpha1().HCPPolicies("hcp").Get(context.TODO(), policy_name, metav1.GetOptions{})
	if err != nil {
		klog.Info(err)
	}
	return policy
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
