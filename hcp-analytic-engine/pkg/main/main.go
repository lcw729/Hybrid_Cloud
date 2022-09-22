package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/KETI-Hybrid/hcp-analytic-engine-v1/pkg/handler"
	"github.com/KETI-Hybrid/hcp-analytic-engine-v1/pkg/metric"

	resource "github.com/KETI-Hybrid/hcp-pkg/hcp-resource/hcppolicy"

	clusterManager "github.com/KETI-Hybrid/hcp-pkg/util/clusterManager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var cm, _ = clusterManager.NewClusterManager()

func main() {
	// var jsonarray PodMetric
	// var cluster_list = []string{"eks-cluster1"}
	// cluster_list 생성 우선 gke-cluster, aks-cluster, eks-cluster 가 저장되어있다고 가정

	// kubeconfig := os.Getenv("KUBECONFIG")
	// config, err := rest.InClusterConfig()

	// creates the clientset
	clientset := cm.Host_kubeClient

	// podsInNode, _ := hostKubeClient.CoreV1().Pods("").List(metav1.ListOptions{})
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
		LabelSelector: "metadata.labels.app=name=hcp-metric-collector",
	})
	for _, p := range pods.Items {
		klog.Info(p.GetName())
	}
	if err != nil {
		panic(err.Error())
	}
	// node 리스트 출력
	// nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	// for _, n := range nodes.Items {
	// 	klog.Info(n.GetName())
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
		// 	klog.Error(err)
		// } else {
		// 	klog.Info(string(output))
		// }
		// clusterList, err := ioutil.ReadFile("usr/local/bin/cluster_list1.txt")
		// clusterList, err := ioutil.ReadFile("~/../usr/local/bin/kubectl_version.txt")
		// if err != nil {
		// 	panic(err)
		// }
		// klog.Info(string(clusterList))
		// for i := 0; i < len(cluster_list); i++ {

		// 	// Calculate_WatchingLevel(podNum[i], cluster_list[i])
		// }

		//노드 메트릭
		// Calculate_Node_Metric(cluster_list[0])

		// Calculate_Cluster_Metric(cluster_list[0])
		// 클러스터 메트릭 수집
		// Calculate_Cluster_Metric(cluster_list[0])
		// klog.Info(cluster_list[0])
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
	// 			klog.Info(p.GetName())
	// 		}
	// 		if err != nil {
	// 			panic(err.Error())
	// 		}
	// 	}
	// }
	// klog.Info(GetClusterExpandingCrierion())
	// klog.Info(GetDefaultNodeOption())

}

func GetClusterExpandingCrierion() (float32, float32) {
	var max_cpu float32
	var max_memory float32
	var extra float32
	var cluster_cpu_criterion float32
	var cluster_mem_criterion float32

	// get initial-setting policy
	initial_setting, err := resource.GetHCPPolicy(*cm.HCPPolicy_Client, "initial-setting")
	if err != nil {
		klog.Info(nil)
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
		klog.Info(cluster_cpu_criterion)
		cluster_mem_criterion = convertTOKiB(max_memory) * extra / 100
		klog.Info(cluster_mem_criterion)

		return cluster_cpu_criterion, cluster_mem_criterion
	}
	return -1, -1
}

func GetDefaultNodeOption() (int, int, int) {
	var option_cpu, option_mem, option_podnum int
	// get initial-setting policy
	initial_setting, err := resource.GetHCPPolicy(*cm.HCPPolicy_Client, "initial-setting")
	if err != nil {
		klog.Info(nil)
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
	klog.Info(option)
	node_option, err := resource.GetHCPPolicy(*cm.HCPPolicy_Client, "node-option")
	if err != nil {
		klog.Info(nil)
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

func Calculate_Node_Metric(cluster_name string) {
	var jsonarray metric.NodeMetric
	nodeNum := 2
	jsonByteArray := handler.GetResource(1, "eks-cluster", "nodes")
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		panic(err)
	}

	for i := 0; i < nodeNum; i++ {

		klog.Info("------------------------------------------------------------------")
		klog.Info("[", i+1, "] Node:", jsonarray.Nodemetrics[i].Node, " Metric Information")
		klog.Info("NodeMetric        :", jsonarray.Nodemetrics[i].Node)
		klog.Info("Time              :", jsonarray.Nodemetrics[i].Time)
		klog.Info("Cpu               :", jsonarray.Nodemetrics[i].CPU)
		klog.Info("Cluster           :", jsonarray.Nodemetrics[i].Cluster)
		klog.Info("Memory           :", jsonarray.Nodemetrics[i].Memory)

		// memoryusage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		// if err != nil {
		// 	// handle error
		// 	klog.Error(err)
		// 	os.Exit(2)
		// }

		// // memoryusage, err = strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryUsageBytes, "Mi"), 32)
		// // if err != nil {
		// // 	// handle error
		// // 	klog.Error(err)
		// // 	os.Exit(2)
		// // }
		// memorytotal, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		// if err != nil {
		// 	// handle error
		// 	klog.Error(err)
		// 	os.Exit(2)
		// }

		// memorytotal, err = strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryAvailableBytes, "Mi"), 32)
		// if err != nil {
		// 	// handle error
		// 	klog.Error(err)
		// 	os.Exit(2)
		// }

		CpuUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].CPU.CPUUsageNanoCores, "n"), 32)
		if err != nil {
			// handle error
			klog.Error(err)
			os.Exit(2)
		}
		MemoryUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		if err != nil {
			// handle error
			klog.Error(err)
			os.Exit(2)
		}
		klog.Info("MemoryUsage: ", MemoryUsage)
		MemoryWorkingSetBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryWorkingSetBytes, "KiMi"), 32)

		if err != nil {
			// handle error
			klog.Error(err)
			os.Exit(2)
		}
		klog.Info("MemoryWorkingSetBytes: ", MemoryWorkingSetBytes)
		MemoryAvailableBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		if err != nil {
			// handle error
			klog.Error(err)
			os.Exit(2)
		}
		klog.Info("MemoryAvailableBytes: ", MemoryAvailableBytes)
		CpuUsage = CpuUsage / 1000000000
		MemoryUsage = MemoryUsage / (MemoryWorkingSetBytes + MemoryAvailableBytes)
		klog.Info("-------------CpuUsage:", CpuUsage, "---------------")
		klog.Info("-------------MemoryUsage:", MemoryUsage, "---------------")

	}
}

func Calculate_Cluster_Metric(cluster_name string) {
	var jsonarray metric.NodeMetric
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
		klog.Info("ClusterMemoryTotal: ", ClusterMemoryTotal)

		MemoryUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		klog.Info("MemoryUsage: ", MemoryUsage)
		if err != nil {
			// handle error
			klog.Error(err)
			os.Exit(2)
		}

		MemoryWorkingSetBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryWorkingSetBytes, "KiMi"), 32)
		klog.Info("MemoryWorkingSetBytes: ", MemoryWorkingSetBytes)
		if err != nil {
			// handle error
			klog.Error(err)
			os.Exit(2)
		}

		MemoryAvailableBytes, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Nodemetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		klog.Info("MemoryAvailableBytes: ", MemoryAvailableBytes)
		if err != nil {
			// handle error
			klog.Error(err)
			os.Exit(2)
		}

		ClusterMemoryTotal = MemoryWorkingSetBytes + MemoryAvailableBytes + ClusterMemoryTotal
		klog.Info("ClusterMemoryTotal: ", ClusterMemoryTotal)
		ClusterMemoryUsage = ClusterMemoryUsage + MemoryUsage

	}

	Usage = ClusterMemoryUsage / ClusterMemoryTotal

	klog.Info("-------------------Cluster:", cluster_name, "Metric information-------------------")
	klog.Info("Print Cluster : ", cluster_name)
	klog.Info("Cluster total Memory            : ", ClusterMemoryTotal)
	klog.Info("Cluster USED Memory            : ", ClusterMemoryUsage)
	klog.Info("Cluster", cluster_name, "Memory Usage: ", Usage)
}

func Calculate_WatchingLevel(podNum int, clusterName string) {
	var watchingLevel = 0
	var jsonarray metric.PodMetric
	jsonByteArray := handler.GetResource(podNum, clusterName, "pods")
	stringArray := string(jsonByteArray[:])
	if err := json.Unmarshal([]byte(stringArray), &jsonarray); err != nil {
		panic(err)
	}
	klog.Info("----------------------------------------------------------")
	klog.Info("Print Cluster : ", clusterName)
	klog.Info("")

	//파드 개수를 가져와서 for문의 변수로 넣어주어야 함

	for i := 0; i < podNum; i++ {

		klog.Info("------------------------------------------------------------------")
		klog.Info("[", i+1, "] Pod:", jsonarray.Podmetrics[i].Pod, " Metric Information")
		klog.Info("PodMetric         :", jsonarray.Podmetrics[i].Pod)
		klog.Info("Time              :", jsonarray.Podmetrics[i].Time)
		klog.Info("Cpu               :", jsonarray.Podmetrics[i].CPU)
		klog.Info("Cluster           :", jsonarray.Podmetrics[i].Cluster)
		klog.Info("Memory            :", jsonarray.Podmetrics[i].Memory)
		klog.Info("Namespace         :", jsonarray.Podmetrics[i].Namespace)

		// memoryusage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryUsageBytes, "KiMi"), 32)
		// if err != nil {
		// 	// handle error
		// 	klog.Error(err)
		// 	os.Exit(2)
		// }

		// // memoryusage, err = strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryUsageBytes, "Mi"), 32)
		// // if err != nil {
		// // 	// handle error
		// // 	klog.Error(err)
		// // 	os.Exit(2)
		// // }
		// memorytotal, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryAvailableBytes, "KiMi"), 32)
		// if err != nil {
		// 	// handle error
		// 	klog.Error(err)
		// 	os.Exit(2)
		// }

		// memorytotal, err = strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].Memory.MemoryAvailableBytes, "Mi"), 32)
		// if err != nil {
		// 	// handle error
		// 	klog.Error(err)
		// 	os.Exit(2)
		// }

		CpuUsage, err := strconv.ParseFloat(strings.TrimRight(jsonarray.Podmetrics[i].CPU.CPUUsageNanoCores, "n"), 32)
		if err != nil {
			// handle error
			klog.Error(err)
			os.Exit(2)
		}
		CpuUsage = CpuUsage / 1000000000
		if CpuUsage != 0 {
			klog.Info("CpuUsage          : ", CpuUsage)
			switch {
			case CpuUsage <= 0.2:
				watchingLevel = 1
				klog.Info("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage <= 0.4:
				watchingLevel = 2
				klog.Info("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage <= 0.6:
				watchingLevel = 3
				klog.Info("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage <= 0.8:
				watchingLevel = 4
				klog.Info("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage <= 1:
				watchingLevel = 5
				klog.Info("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
			case CpuUsage > 1:
				watchingLevel = 5
				klog.Info("Pod Name ", jsonarray.Podmetrics[i].Pod, "watching Level is ", watchingLevel)
				klog.Info("Watching Level is over 5")
				klog.Info("!!!!!!!!!!!!!!!", jsonarray.Podmetrics[i].Pod, "!!!!!!!!!!!!!!!!!!!")
			}
		} else {
			klog.Info("This Pod use 0 Cpu nanocores")
		}
	}
	klog.Info("------------------------------------------")
}

func AutoSacling() {
	//디플로이먼트 이름을 넘겨주어야 함
}
