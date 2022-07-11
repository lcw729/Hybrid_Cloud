package main

import (
	"Hybrid_Cloud/hcp-analytic-engine/pkg/handler"
	resource "Hybrid_Cloud/hcp-resource/hcppolicy"
	cobrautil "Hybrid_Cloud/hybridctl/util"
	resourcev1alpha1 "Hybrid_Cloud/pkg/apis/resource/v1alpha1"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
)

func main() {
	// var jsonarray PodMetric
	// var cluster_list = []string{"eks-cluster1"}
	// cluster_list 생성 우선 gke-cluster, aks-cluster, eks-cluster 가 저장되어있다고 가정

	// kubeconfig := os.Getenv("KUBECONFIG")
	// config, err := rest.InClusterConfig()
	config, err := cobrautil.BuildConfigFromFlags("master", "/root/.kube/config")
	// config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset := kubernetes.NewForConfigOrDie(config)

	// podsInNode, _ := hostKubeClient.CoreV1().Pods("").List(metav1.ListOptions{})
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: "metadata.labels.app=name=hcp-metric-collector",
	})
	for _, p := range pods.Items {
		fmt.Println(p.GetName())
	}
	if err != nil {
		panic(err.Error())
	}
	// node 리스트 출력
	// nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	// for _, n := range nodes.Items {
	// 	fmt.Println(n.GetName())
	// }
	// if err != nil {
	// 	panic(err.Error())
	// }

	for {
		// cmd := exec.Command("kubectl", "config", "get-contexts", "--output='name'", ">", "cluster_list.txt")
		// cmd := exec.Command("kubectl", "version", ">", "kubectl_version.txt")
		// cmd.Dir = "usr/local/bin"
		// output, err := cmd.Output()
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(string(output))
		// }
		// clusterList, err := ioutil.ReadFile("usr/local/bin/cluster_list1.txt")
		// clusterList, err := ioutil.ReadFile("~/../usr/local/bin/kubectl_version.txt")
		// if err != nil {
		// 	panic(err)
		// }
		// fmt.Println(string(clusterList))
		// for i := 0; i < len(cluster_list); i++ {

		// 	// Calculate_WatchingLevel(podNum[i], cluster_list[i])
		// }

		//노드 메트릭
		// Calculate_Node_Metric(cluster_list[0])

		// Calculate_Cluster_Metric(cluster_list[0])
		// 클러스터 메트릭 수집
		// Calculate_Cluster_Metric(cluster_list[0])
		// fmt.Println(cluster_list[0])
		time.Sleep(2 * time.Second)
	}

	// lcw

	// cm, err := clusterManager.NewClusterManager()
	// if err != nil {
	// 	klog.Error(err)
	// } else {
	// 	clusters := cm.Cluster_list
	// 	for _, cluster := range clusters.Items {
	// 		config, err := cobrautil.BuildConfigFromFlags(cluster.ObjectMeta.Name, "/root/.kube/config")
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
		fmt.Println("Cluster           :", jsonarray.Nodemetrics[i].Cluster)
		fmt.Println("Memory           :", jsonarray.Nodemetrics[i].Memory)

		// memoryusage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		// if err != nil {
		// 	// handle error
		// 	fmt.Println(err)
		// 	os.Exit(2)
		// }

		// // memoryusage, err = strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryUsageBytes, "Mi"), 32)
		// // if err != nil {
		// // 	// handle error
		// // 	fmt.Println(err)
		// // 	os.Exit(2)
		// // }
		// memorytotal, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		// if err != nil {
		// 	// handle error
		// 	fmt.Println(err)
		// 	os.Exit(2)
		// }

		// memorytotal, err = strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryAvailableBytes, "Mi"), 32)
		// if err != nil {
		// 	// handle error
		// 	fmt.Println(err)
		// 	os.Exit(2)
		// }

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
		fmt.Println("MemoryUsage: ", MemoryUsage)
		MemoryWorkingSetBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryWorkingSetBytes, "KiMi"), 32)

		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("MemoryWorkingSetBytes: ", MemoryWorkingSetBytes)
		MemoryAvailableBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		fmt.Println("MemoryAvailableBytes: ", MemoryAvailableBytes)
		CpuUsage = CpuUsage / 1000000000
		MemoryUsage = MemoryUsage / (MemoryWorkingSetBytes + MemoryAvailableBytes)
		fmt.Println("-------------CpuUsage:", CpuUsage, "---------------")
		fmt.Println("-------------MemoryUsage:", MemoryUsage, "---------------")

	}
}

func Calculate_Cluster_Metric(cluster_name string) {
	var jsonarray NodeMetric
	node_num := 2
	jsonByteArray := handler.GetResource(1, "eks-cluster", "nodes")
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		panic(err)
	}

	ClusterMemoryTotal := 0.0
	ClusterMemoryUsage := 0.0
	var Usage float64
	for i := 0; i < node_num; i++ {
		fmt.Println("ClusterMemoryTotal: ", ClusterMemoryTotal)

		MemoryUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		fmt.Println("MemoryUsage: ", MemoryUsage)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}

		MemoryWorkingSetBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryWorkingSetBytes, "KiMi"), 32)
		fmt.Println("MemoryWorkingSetBytes: ", MemoryWorkingSetBytes)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}

		MemoryAvailableBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		fmt.Println("MemoryAvailableBytes: ", MemoryAvailableBytes)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}

		ClusterMemoryTotal = MemoryWorkingSetBytes + MemoryAvailableBytes + ClusterMemoryTotal
		fmt.Println("ClusterMemoryTotal: ", ClusterMemoryTotal)
		ClusterMemoryUsage = ClusterMemoryUsage + MemoryUsage

	}

	Usage = ClusterMemoryUsage / ClusterMemoryTotal

	fmt.Println("-------------------Cluster:", cluster_name, "Metric information-------------------")
	fmt.Println("Print Cluster : ", cluster_name)
	fmt.Println("Cluster total Memory            : ", ClusterMemoryTotal)
	fmt.Println("Cluster USED Memory            : ", ClusterMemoryUsage)
	fmt.Println("Cluster", cluster_name, "Memory Usage: ", Usage)
}

func Calculate_WatchingLevel(podNum int, clusterName string) {
	var watchingLevel = 0
	var jsonarray PodMetric
	jsonByteArray := handler.GetResource(podNum, clusterName, "pods")
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		panic(err)
	}
	fmt.Println("----------------------------------------------------------")
	fmt.Println("Print Cluster : ", clusterName)
	fmt.Println("")

	//파드 개수를 가져와서 for문의 변수로 넣어주어야 함

	for i := 0; i < podNum; i++ {

		fmt.Println("------------------------------------------------------------------")
		fmt.Println("[", i+1, "] Pod:", jsonarray.Podmetrics[i].Pod, " Metric Information")
		fmt.Println("PodMetric         :", jsonarray.Podmetrics[i].Pod)
		fmt.Println("Time              :", jsonarray.Podmetrics[i].Time)
		fmt.Println("Cpu               :", jsonarray.Podmetrics[i].CPU)
		fmt.Println("Cluster           :", jsonarray.Podmetrics[i].Cluster)
		fmt.Println("Memory            :", jsonarray.Podmetrics[i].Memory)
		fmt.Println("Namespace         :", jsonarray.Podmetrics[i].Namespace)

		// memoryusage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		// if err != nil {
		// 	// handle error
		// 	fmt.Println(err)
		// 	os.Exit(2)
		// }

		// // memoryusage, err = strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryUsageBytes, "Mi"), 32)
		// // if err != nil {
		// // 	// handle error
		// // 	fmt.Println(err)
		// // 	os.Exit(2)
		// // }
		// memorytotal, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		// if err != nil {
		// 	// handle error
		// 	fmt.Println(err)
		// 	os.Exit(2)
		// }

		// memorytotal, err = strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryAvailableBytes, "Mi"), 32)
		// if err != nil {
		// 	// handle error
		// 	fmt.Println(err)
		// 	os.Exit(2)
		// }

		CpuUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].CPU.CPUUsageNanoCores, "n"), 32)
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		CpuUsage = CpuUsage / 1000000000
		if CpuUsage != 0 {
			fmt.Println("CpuUsage          : ", CpuUsage)
			switch {
			case CpuUsage <= 0.2:
				watchingLevel = 1
				fmt.Println("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage <= 0.4:
				watchingLevel = 2
				fmt.Println("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage <= 0.6:
				watchingLevel = 3
				fmt.Println("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage <= 0.8:
				watchingLevel = 4
				fmt.Println("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage <= 1:
				watchingLevel = 5
				fmt.Println("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage > 1:
				watchingLevel = 5
				fmt.Println("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
				fmt.Println("Watching Level is over 5")
				fmt.Println("!!!!!!!!!!!!!!!", jsonarray.Podmetrics[i].Pod, "!!!!!!!!!!!!!!!!!!!")
			}
		} else {
			fmt.Println("This Pod use 0 Cpu nanocores")
		}
	}
	fmt.Println("------------------------------------------")
}

func AutoSacling() {
	//디플로이먼트 이름을 넘겨주어야 함
}
