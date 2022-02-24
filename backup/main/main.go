package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// "Hybrid_Cloud/analytic-engine/analyticEngine"
	// "Hybrid_Cloud/hcp-scheduler/pkg/policy"
	"Hybrid_Cloud/hcp-scheduler/backup/policy"
	scheduler "Hybrid_Cloud/hcp-scheduler/backup/scheduler"
	algopb "Hybrid_Cloud/protos/v1/algo"
	"Hybrid_Cloud/resource"

	"google.golang.org/grpc"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// // grpc "google.golang.org/grpc"
	// "context"
	// "fmt"
	// "log"
	// "time"
	// v1 "k8s.io/api/apps/v1"
	// "k8s.io/client-go/kubernetes"
)

func main() {
	replicas := int32(3)

	resourceList := make(corev1.ResourceList)
	resourceList.Cpu().Set(2)
	resourceList.Memory().Set(2)
	deployment := v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-deployment",
			Namespace: "default",
			Labels:    map[string]string{"app": "nginx"},
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": "nginx"},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:1.14.2",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits:   resourceList,
								Requests: resourceList,
							},
						},
					},
					// NodeName: "cluster2-worker2",
				},
			},
		},
	}
	// resource.UpdateDeployment("kube-master", &deployment, replicas)
	err := resource.CreateDeployment("kube-master", "cluster2-worker2", &deployment)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// p, err := resource.GetPod("kube-master", "nginx-deployment-69f8d49b75-qhrtd", "default")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// d, err := resource.GetDeployment("kube-master", p)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	resource.UpdateDeployment("kube-master", d, 5)
	// }
	// err = resource.CreateDeployment("aks-master", "aks-agentpool-21474300-vmss000003", d)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// createDeployment("kube-master")
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := algopb.NewAlgoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	optimalArrangement(c, ctx)
	// ResourceExtensionSchedule(c, ctx)
	// ResourceConfigurationSchedule(c, ctx)

}

func ResourceExtensionSchedule(c algopb.AlgoClient, ctx context.Context) {
	fmt.Println("[ Scheduler Start ]")
	fmt.Println("[step 1] Check Policy from Policy manager - calculation cycle")
	cycle := policy.GetCycle()
	fmt.Println(cycle)
	if cycle > 0 {
		fmt.Println("-------------------------LOOP START----------------------------")
		for {
			time.Sleep(time.Second * time.Duration(cycle))
			fmt.Println("[step 2] Get Cluster WeightResult")
			// 가중치 계산 결과
			var in = algopb.ClusterWeightCalculatorRequest{}
			r, err := c.ClusterWeightCalculator(ctx, &in)
			if err != nil {
				log.Fatalf("could not request: %v", err)
			}

			log.Printf("Config: %v", r.GetWeightResult())
			scheduler.Resourcebalancingcontroller()
			fmt.Println("---------------------------------------------------------------")
		}
	}
}

func optimalArrangement(c algopb.AlgoClient, ctx context.Context) {
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[ Scheduler Start ]")
	var in = algopb.OptimalArrangementRequest{}
	r, err := c.OptimalArrangement(ctx, &in)
	if err != nil {
		log.Fatalf("could not request: %v", err)
	}

	log.Printf("Config: %v ", r)
}
