package main

import (
	// "Hybrid_Cluster/analytic-engine/analyticEngine"

	"Hybrid_Cluster/hcp-scheduler/pkg/policy"
	scheduler "Hybrid_Cluster/hcp-scheduler/pkg/scheduler"
	algopb "Hybrid_Cluster/protos/v1/algo"

	grpc "google.golang.org/grpc"

	"context"
	"fmt"
	"log"
	"time"
)

func main() {

	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := algopb.NewAlgoClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	ResourceExtensionSchedule(c, ctx)
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

// func ResourceConfigurationSchedule(c algopb.AlgoClient, ctx context.Context) {
// 	fmt.Println("-----------------------------------------")
// 	fmt.Println("[step 1] Get MultiMetric")
// 	// monitoringEngine.MetricCollector()
// 	fmt.Println("-----------------------------------------")
// 	fmt.Println("[step 2] Get Policy - watching level & warning level")
// 	fmt.Println("< Watching Level > \n", policy.GetWatchingLevel())
// 	fmt.Println("< Warning  Level > \n", policy.GetWarningLevel())
// 	fmt.Println("[step 3] Calculate watching level")
// 	var targetCluster []*algopb.ClusterInfo

// 	targetCluster = append(targetCluster, &algopb.ClusterInfo{
// 		ClusterId: 1,
// 		ClusterName :"cluster1",
// 	})
// 	targetCluster = append(targetCluster, &algopb.ClusterInfo{
// 		ClusterId: 2,
// 		ClusterName :"cluster2",
// 	})
// 	targetCluster = append(targetCluster, &algopb.ClusterInfo{
// 		ClusterId: 3,
// 		ClusterName :"cluster3",
// 	})
// 	fmt.Println(targetCluster)
// 	in := &algopb.ResourceConfigurationManagementRequest{
// 		TargetCluster: targetCluster,
// 	}
// 	r, err := c.ResourceConfigurationManagement(ctx, in)
// 	if err != nil {
// 		log.Fatalf("could not request: %v", err)
// 	}

// 	log.Printf("Status: %v", r.Status)
// 	log.Printf("TestMessage: %v", r.TestMessage)
// }

func optimalArrangement() {
	fmt.Println("---------------------------------------------------------------")
	fmt.Println("[ Scheduler Start ]")
	fmt.Println("[step 1] Check Policy from Policy manager - DRF, Affinity")
	fmt.Println("----> Policy Value: Affinity")
	fmt.Println("----> Policy Value: DRF")
	fmt.Println("[step 2] Call Affinity Calulator")
	// analyticEngine.AffinityCalculator()
	fmt.Println("[step 3] Profiling Pod & Node")
	// analyticEngine.DRF()
	fmt.Println("[step 4] Checking Pending POD Queue")
	fmt.Println("----> [case 1] If there are no suitable resources, wait ")
	fmt.Println("----> [case 2] Select suitable target resources and resource request ")
}
