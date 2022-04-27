package main

import (
	resource "Hybrid_Cloud/hcp-analytic-engine/pkg/autoscaler"
	"Hybrid_Cloud/hcp-analytic-engine/pkg/backup/algorithm"
	cm "Hybrid_Cloud/util/clusterManager"
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 3] Calculate resource weight")
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 4] Send weight calculation result to Scheduler (Resource Balancing Controller)")
	fmt.Println("--Resource Weight Result--")
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
		fmt.Println(c.ClusterInfo, n.Score)
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

		fmt.Println(extra)
		cpu = cpu * (100 - extra) / 100
		mem = mem * (100 - extra) / 100
		fmt.Println(cpu, mem)
	*/

	// HPA/VPA 함수 사용 예시
	var pod *v1.Pod
	cluster := "aks-master"
	test_pod_name := "nginx-deploy-6d4c4cc4b8-98zrr"
	ns := "default"

	clustermanager, err := cm.NewClusterManager()
	clientset := clustermanager.Cluster_kubeClients[cluster]
	deployment, _ := clientset.AppsV1().Pods(ns).Get(context.TODO(), test_pod_name, metav1.GetOptions{})

	if err != nil {
		fmt.Println(err)
	} else {
		autoscaler := resource.NewAutoScaler(cluster, pod, ns)
		resource.AutoscalerMap[cluster] = autoscaler

		if bol, _ := algorithm.WatchingLevelCalculator(); bol {
			autoscaler.WarningCountPlusOne(pod)
			autoscaler.AutoScaling(pod)
		}
	}

	/*
		lis, err := net.Listen("tcp", ":"+portNumber)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()
		algopb.RegisterAlgoServer(grpcServer, &algoServer{})

		log.Printf("start gRPC server on %s port", portNumber)
		fmt.Println("[step 1] Get ResourceConfigurationCycle Policy")
		cycle := policy.GetCycle()
		if cycle > 0 {
			for {
				time.Sleep(time.Second * time.Duration(cycle))
				fmt.Println("-------------------------LOOP START----------------------------")
				algorithm.WatchingLevelCalculator()
			}
		} else {
			fmt.Println("Error : Cycle should be positive")
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	*/
}
