package main

import (
	// algopb "Hybrid_Cluster/protos/v1/algo"
	resource "Hybrid_Cluster/hcp-scheduler/pkg/resource"
	resourcev1alpha1 "Hybrid_Cluster/pkg/apis/resource/v1alpha1"
	hasv1alpha1 "Hybrid_Cluster/pkg/client/resource/v1alpha1/clientset/versioned"
	cm "Hybrid_Cluster/util/clusterManager"
	"context"
	"fmt"

	autoscaling "k8s.io/api/autoscaling/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	vpav1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"

	hpav1 "k8s.io/api/autoscaling/v1"
)

/*
const portNumber = "9000"

type algoServer struct {
	algopb.AlgoServer
}


// 리소스 확장 기술 -- 가중치 계산 [가중치 계산 결과 넘겨줌]
// scheduler -> analytic Engine
func (a *algoServer) ClusterWeightCalculator(ctx context.Context, in *algopb.ClusterWeightCalculatorRequest) (*algopb.ClusterWeightCalculatorResponse, error) {
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[step 2] Get MultiMetric")
	monitoringEngine.MetricCollector()
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
	cm := cm.NewClusterManager()
	master_config := cm.Host_config
	// client := kubernetes.NewForConfigOrDie(master_config)
	cluster := "kube-master"
	pod := "nginx-deployment-69f8d49b75-hc7pd"
	ns := "default"
	p, err := resource.GetPod(cluster, pod, ns)
	if err != nil {
		fmt.Println(err)
	}
	d, err := resource.GetDeployment(cluster, p)
	if err != nil {
		fmt.Println(err)
	}
	deploymentName := "nginx-deployment"
	updateMode := vpav1.UpdateModeAuto
	vpa := vpav1.VerticalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "autoscaling.k8s.io/v1",
			Kind:       "VerticalPodAutoscaler",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: vpav1.VerticalPodAutoscalerSpec{
			TargetRef: &autoscaling.CrossVersionObjectReference{
				APIVersion: d.APIVersion,
				Kind:       d.Kind,
				Name:       d.Name,
			},
			UpdatePolicy: &vpav1.PodUpdatePolicy{
				UpdateMode: &updateMode,
			},
		},
	}
	instance := &resourcev1alpha1.HCPHybridAutoScaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment" + "-vpa",
		},
		Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
			WarningCount: 3,
			CurrentStep:  "HAS", // HAS -> Sync -> Done
			ScalingOptions: resourcev1alpha1.ScalingOptions{
				VpaTemplate: vpa,
			},
		},
	}
	hasv1alpha1clientset, err := hasv1alpha1.NewForConfig(master_config)
	if err != nil {
		fmt.Println(err)
	}
	newhas, err := hasv1alpha1clientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Create(context.TODO(), instance, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("create %s Done\n", newhas.Name)
	}

}

func hpatest() {
	cm := cm.NewClusterManager()
	master_config := cm.Host_config
	// client := kubernetes.NewForConfigOrDie(master_config)

	var num int32 = 1
	var min = &num
	hpa := hpav1.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-deployment",
		},
		Spec: hpav1.HorizontalPodAutoscalerSpec{
			MinReplicas: min,
			MaxReplicas: 10,
		},
	}
	hpa.Spec.ScaleTargetRef.APIVersion = "apps/v1"
	hpa.Spec.ScaleTargetRef.Kind = "Deployment"
	hpa.Spec.ScaleTargetRef.Name = "nginx-deployment"
	instance := &resourcev1alpha1.HCPHybridAutoScaler{
		ObjectMeta: metav1.ObjectMeta{
			Name: hpa.Spec.ScaleTargetRef.Name,
		},
		Spec: resourcev1alpha1.HCPHybridAutoScalerSpec{
			WarningCount: 1,
			CurrentStep:  "HAS", // HAS -> Sync -> Done
			ScalingOptions: resourcev1alpha1.ScalingOptions{
				HpaTemplate: hpa,
			},
		},
	}
	hasv1alpha1clientset, err := hasv1alpha1.NewForConfig(master_config)
	if err != nil {
		fmt.Println(err)
	}
	newhas, err := hasv1alpha1clientset.HcpV1alpha1().HCPHybridAutoScalers("hcp").Create(context.TODO(), instance, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("create %s Done\n", newhas.Name)
	}

}
