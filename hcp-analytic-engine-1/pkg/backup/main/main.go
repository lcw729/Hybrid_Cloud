package main

import (
	resourcev1alpha1 "Hybrid_Cloud/pkg/apis/resource/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
)

// policy "Hybrid_Cloud/hcp-analytic-engine/pkg/policy"
//"Hybrid_Cloud/hcp-analytic-engine/pkg/autoscaling"
//algopb "Hybrid_Cloud/protos/v1/algo"

/*
const portNumber = "9000"

type algoServer struct {
	algopb.AlgoServer
}
*/

/*
// 리소스 확장 기술 -- 가중치 계산 [가중치 계산 결과 넘겨줌]
// scheduler -> analytic Engine
func (a *algoServer) ClusterWeightCalculator(ctx context.Context, in *algopb.ClusterWeightCalculatorRequest) (*algopb.ClusterWeightCalculatorResponse, error) {
	klog.Info("---------------------------------------------------------------")
	klog.Info("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	klog.Info("---------------------------------------------------------------")
	klog.Info("[step 3] Calculate resource weight")
	klog.Info("---------------------------------------------------------------")
	klog.Info("[step 4] Send weight calculation result to Scheduler (Resource Balancing Controller)")
	klog.Info("--Resource Weight Result--")
	weightResult := make([]*algopb.WeightResult, 4)
	weightResult[0] = &algopb.WeightResult{
		ClusterId:     1,
		ClusterName:   "cluster1",
		ClusterWeight: 30,
	}
	weightResult[1] = &algopb.WeightResult{
		ClusterId:     2,
		ClusterName:   "cluster2",
		ClusterWeight: 20,
	}
	weightResult[2] = &algopb.WeightResult{
		ClusterId:     3,
		ClusterName:   "cluster3",
		ClusterWeight: 25,
	}
	weightResult[3] = &algopb.WeightResult{
		ClusterId:     4,
		ClusterName:   "cluster4",
		ClusterWeight: 25,
	}

	return &algopb.ClusterWeightCalculatorResponse{
		WeightResult: weightResult,
	}, nil
}

func (a *algoServer) OptimalArrangement(ctx context.Context, in *algopb.OptimalArrangementRequest) (*algopb.OptimalArrangementResponse, error) {
	var c *util.Cluster
	var n *util.NodeScore
	if algorithm.OptimalArrangementAlgorithm() {
		c, n = algorithm.OptimalNodeSelector()
		klog.Info(c.ClusterInfo, n.Score)
	}
	return &algopb.OptimalArrangementResponse{
		Status: true,
		Cluster: &algopb.Cluster{
			ClusterInfo: (*algopb.ClusterInfo)(c.ClusterInfo),
		},
		Node: &algopb.NodeScore{
			NodeId: n.NodeId,
			Score:  n.Score,
		},
	}, nil
}
*/

func main() {

	//algorithm.Affinity()
	/*
		cpu, _ := policy.GetInitialSettingValue("max_cpu")
		mem, _ := policy.GetInitialSettingValue("max_memory")
		extra, _ := policy.GetInitialSettingValue("extra")

		klog.Info(extra)
		cpu = cpu * (100 - extra) / 100
		mem = mem * (100 - extra) / 100
		klog.Info(cpu, mem)
	*/

	// HPA/VPA 함수 사용 예시
	// cluster := "aks-master"
	// test_dep_name := "nginx-deploy"
	// ns := "default"

	//clustermanager, err := cm.NewClusterManager()
	//clientset := clustermanager.Cluster_kubeClients[cluster]
	//deployment, _ := clientset.AppsV1().Deployments(ns).Get(context.TODO(), test_dep_name, metav1.GetOptions{})

	// var jsonarray PodMetric
	// var cluster_list = []string{"gke-cluster"}
	// cluster_list 생성 우선 gke-cluster, aks-cluster, eks-cluster 가 저장되어있다고 가정
	// var podNum = 2
	/*
		for {
			// cmd := exec.Command("kubectl", "config", "get-contexts", "--output='name'", ">", "cluster_list.txt")
			cmd := exec.Command("kubectl", "version", ">", "kubectl_version.txt")
			cmd.Dir = "usr/local/bin"
			output, err := cmd.Output()
			if err != nil {
				klog.Error(err)
			} else {
				klog.Info(string(output))
			}
	*/
	// clusterList, err := ioutil.ReadFile("usr/local/bin/cluster_list1.txt")
	// clusterList, err := ioutil.ReadFile("~/../usr/local/bin/kubectl_version.txt")
	// if err != nil {
	// 	panic(err)
	// }
	// klog.Info(string(clusterList))

	//for i := 0; i < len(cluster_list); i++ {
	// bol, pod, namespace, _ := algorithm.Calculate_WatchingLevel(podNum[i], cluster_list[i])
	// klog.Info(bol, pod, namespace)

	// deployment := hcpdeploymentToDeployment(hcpdeployment)

	/*
			if bol, pod, namespace, _ := algorithm.Calculate_WatchingLevel(j, podNum, cluster_list[i]); bol {
				// 1. autoscalerMap에 cluster 등록되어있는지 확인

					po, _ := kuberesourcepo.GetPod(cluster_list[i], pod, namespace)
					deployment, err := kuberesourcedeploy.GetDeployment(cluster_list[i], po)

				if err == nil {
		if resource.AutoscalerMap[cluster_list[i]] == nil {
			klog.Info("===========no autoscaler===========")
			// autoscalerMap에 cluster autoscaler 저장
			autoscaler := resource.NewAutoScaler()
			autoscaler.RegisterDeploymentToAutoScaler(&deployment)
			resource.AutoscalerMap[cluster_list[i]] = autoscaler
			autoscaler.WarningCountPlusOne(&deployment)
			autoscaler.AutoScaling(&deployment)
			klog.Info("current warningcount is ", resource.AutoscalerMap[cluster_list[i]].GetWarningCount(&deployment))
			klog.Info("===================================")
		} else {
			autoscaler := resource.AutoscalerMap[cluster_list[i]]
	*/

	/*
		} else {
					klog.Error(err)
				}
			}
		}
		time.Sleep(10 * time.Second)
		klog.Info("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
			}
		}
	*/
	/*
		lis, err := net.Listen("tcp", ":"+portNumber)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		algopb.RegisterAlgoServer(grpcServer, &algoServer{})

		log.Printf("start gRPC server on %s port", portNumber)
		klog.Info("[step 1] Get ResourceConfigurationCycle Policy")
		cycle := policy.GetCycle()
		if cycle > 0 {
			for {
				time.Sleep(time.Second * time.Duration(cycle))
				klog.Info("-------------------------LOOP START----------------------------")
				algorithm.WatchingLevelCalculator()
			}
		} else {
			klog.Info("Error : Cycle should be positive")
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	*/
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
