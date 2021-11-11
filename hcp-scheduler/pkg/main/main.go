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
