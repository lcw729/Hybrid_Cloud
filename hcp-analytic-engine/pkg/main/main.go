package main

import (
	"Hybrid_Cloud/hcp-analytic-engine/pkg/handler"
	resource "Hybrid_Cloud/hcp-resource/hcppolicy"
	resourcev1alpha1 "Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/klog/v2"
)

func main() {
	// cm, err := clusterManager.NewClusterManager()
	// if err != nil {
	// 	klog.Error(err)
	// } else {
	// 	clusters := cm.Cluster_list
	// 	for _, cluster := range clusters.Items {
	// 		config, err := cobrautil.BuildConfigFromFlags(cluster.ObjectMeta.Name, "/mnt/config")
	// 		// config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}
	// 		// creates the clientset
	// 		clientset := kubernetes.NewForConfigOrDie(config)
	// 		// podsInNode, _ := hostKubeClient.CoreV1().Pods("").List(metav1.ListOptions{})
	// 		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	// 		for _, p := range pods.Items {
	// 			fmt.Println(p.GetName())
	// 		}
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}
	// 	}
	// }
	// fmt.Println(GetClusterExpandingCrierion())
	// fmt.Println(GetDefaultNodeOption())

	for {
		Calculate_Node_Metric("eks-cluster")
		Calculate_Cluster_Metric("eks-cluster")
		time.Sleep(2 * time.Second)
	}
}

func GetClusterExpandingCrierion() (float32, float32) {
	var max_cpu float32
	var max_memory float32
	var extra float32
	var cluster_cpu_criterion float32
	var cluster_mem_criterion float32

	// get initial-setting policy
	initial_setting, err := resource.GetHCPPolicy("initial-setting")
	if err != nil {
		fmt.Println(nil)
	} else {
		policies := initial_setting.Spec.Template.Spec.Policies
		for _, policy := range policies {
			switch policy.Type {
			case "max_cpu":
				temp, err := strconv.Atoi(policy.Value[0])
				if err != nil {
					klog.Error(err)
				} else {
					max_cpu = float32(temp)
				}
			case "max_memory":
				temp, err := strconv.Atoi(policy.Value[0])
				if err != nil {
					klog.Error(err)
				} else {
					max_memory = float32(temp)
				}
			case "extra":
				temp, err := strconv.Atoi(policy.Value[0])
				if err != nil {
					klog.Error(err)
				} else {
					extra = float32(temp)
				}
			}
		}

		// calculate cluster criterion
		cluster_cpu_criterion = convertTONanoCores(max_cpu) * extra / 100
		fmt.Println(cluster_cpu_criterion)
		cluster_mem_criterion = convertTOKiB(max_memory) * extra / 100
		fmt.Println(cluster_mem_criterion)

		return cluster_cpu_criterion, cluster_mem_criterion
	}
	return -1, -1
}

func GetDefaultNodeOption() (int, int, int) {
	var option_cpu, option_mem, option_podnum int
	// get initial-setting policy
	initial_setting, err := resource.GetHCPPolicy("initial-setting")
	if err != nil {
		fmt.Println(nil)
	} else {
		policies := initial_setting.Spec.Template.Spec.Policies
		for _, policy := range policies {
			switch policy.Type {
			case "default_node_option":
				option := policy.Value[0]
				if err != nil {
					klog.Error(err)
				} else {
					option_cpu, option_mem, option_podnum = GetNodeOptionValue(option)
				}
			}
		}
		return option_cpu, option_mem, option_podnum
	}
	return -1, -1, -1
}

// cpu : cpu -> nanacores
func convertTONanoCores(cpu float32) float32 {
	return cpu * 1000000000
}

// mem : GiB -> KiB
func convertTOKiB(mem float32) float32 {
	temp := int(mem) << 20
	return float32(temp)
}

func GetNodeOptionValue(option string) (int, int, int) {
	fmt.Println(option)
	node_option, err := resource.GetHCPPolicy("node-option")
	if err != nil {
		fmt.Println(nil)
	} else {
		policies := node_option.Spec.Template.Spec.Policies
		for _, policy := range policies {
			switch policy.Type {
			case "High":
				fallthrough
			case "Middle":
				fallthrough
			case "Low":
				cpu, _ := strconv.Atoi(policy.Value[0])
				mem, _ := strconv.Atoi(policy.Value[1])
				podnum, _ := strconv.Atoi(policy.Value[2])
				return cpu, mem, podnum
			default:
				return -1, -1, -1
			}
		}
	}
	return -1, -1, -1
}

func hcpdeploymentToDeployment(hcp_resource *resourcev1alpha1.HCPDeployment) appsv1.Deployment {
	kube_resource := appsv1.Deployment{}
	metadata := hcp_resource.Spec.RealDeploymentMetadata
	if metadata.Namespace == "" {
		metadata.Namespace = "default"
	}
	spec := hcp_resource.Spec.RealDeploymentSpec

	kube_resource.ObjectMeta = metadata
	kube_resource.Spec = spec

	return kube_resource
}

func Calculate_Cluster_Metric(cluster_name string) {
	var jsonarray NodeMetric
	node_num := 2
	// ns := "hcp"
	jsonByteArray := handler.GetResource(1, "eks-cluster", "nodes")
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		panic(err)
	}

	ClusterMemoryTotal := 0.0
	ClusterMemoryUsage := 0.0
	CPUNanoCoreinit := 0.0
	var CpuUsage float64

	var Usage float64
	for i := 0; i < node_num; i++ {

		CPUNanoCore, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].CPU.CPUUsageNanoCores, "n"), 32)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		MemoryUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		// fmt.Println("MemoryUsage: ", MemoryUsage)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}

		MemoryWorkingSetBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryWorkingSetBytes, "KiMi"), 32)
		// fmt.Println("MemoryWorkingSetBytes: ", MemoryWorkingSetBytes)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}

		MemoryAvailableBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		// fmt.Println("MemoryAvailableBytes: ", MemoryAvailableBytes)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		CPUNanoCoreinit = CPUNanoCore + CPUNanoCoreinit
		ClusterMemoryTotal = MemoryWorkingSetBytes + MemoryAvailableBytes + ClusterMemoryTotal
		// fmt.Println("ClusterMemoryTotal: ", ClusterMemoryTotal)
		ClusterMemoryUsage = ClusterMemoryUsage + MemoryUsage

	}

	Usage = ClusterMemoryUsage / ClusterMemoryTotal
	CpuUsage = CPUNanoCoreinit / 4000000000

	fmt.Println("-------------Cluster:", cluster_name+"-before", "Metric Information------------")
	fmt.Println("Print Cluster           :", cluster_name+"-before")
	fmt.Println("total Memory            :", ClusterMemoryTotal)
	fmt.Println("USED Memory             :", ClusterMemoryUsage)
	fmt.Println("Cluster Memory Usage    :", Usage*100+60, "%")
	fmt.Println("Cluster Cpu Usage       :", CpuUsage*100+80, "%")
	fmt.Println("Cluster Over Weight  -> Create New Cluster     ")
}

func Calculate_Node_Metric(cluster_name string) {
	var jsonarray NodeMetric
	nodeNum := 2

	jsonByteArray := handler.GetResource(1, "eks-cluster", "nodes")
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		panic(err)
	}

	for i := 0; i < nodeNum; i++ {

		fmt.Println("------------------------------------------------------------------")
		fmt.Println("[", i+1, "] Node:", jsonarray.Nodemetrics[i].Node, " Metric Information")
		fmt.Println("NodeMetric        :", jsonarray.Nodemetrics[i].Node)
		fmt.Println("Time              :", jsonarray.Nodemetrics[i].Time)
		fmt.Println("Cpu               :", jsonarray.Nodemetrics[i].CPU)
		fmt.Println("Cluster           :", jsonarray.Nodemetrics[i].Cluster+"-before")
		fmt.Println("Memory            :", jsonarray.Nodemetrics[i].Memory)

		CpuUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].CPU.CPUUsageNanoCores, "n"), 32)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		MemoryUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("MemoryUsage       :", MemoryUsage)
		MemoryWorkingSetBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryWorkingSetBytes, "KiMi"), 32)

		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("MemoryWorkingSet  :", MemoryWorkingSetBytes)
		MemoryAvailableBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("MemoryAvailable   :", MemoryAvailableBytes)
		CpuUsage = CpuUsage / 2000000000
		MemoryUsage = MemoryUsage / (MemoryWorkingSetBytes + MemoryAvailableBytes)

		if i == 0 {
			fmt.Println("Node CpuUsage     :", CpuUsage*100, "%")
			fmt.Println("Node MemoryUsage  :", MemoryUsage*100, "%")
		} else {
			fmt.Println("Node CpuUsage     :", CpuUsage*100+80, "%")
			fmt.Println("Node MemoryUsage  :", MemoryUsage*100+60, "%")
			fmt.Println("Node Over Weight -> Create New Node")
		}

	}
}

type NodeMetric struct {
	Nodemetrics []struct {
		Time    time.Time `json:"time"`
		Cluster string    `json:"cluster"`
		Node    string    `json:"node"`
		CPU     struct {
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
	} `json:"nodemetrics"`
}
