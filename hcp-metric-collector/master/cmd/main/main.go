package main

import (
	"Hybrid_Cloud/hcp-metric-collector/master/pkg/metricCollector"
	"Hybrid_Cloud/util/clusterManager"
	"Hybrid_Cloud/util/controller/reshape"
	"fmt"
	"log"
	"os"

	"runtime"

	"admiralty.io/multicluster-controller/pkg/cluster"
	"admiralty.io/multicluster-controller/pkg/manager"
)

const (
	GRPC_PORT = "2051"
)

func main() {
	// logLevel.KetiLogInit()
	// fmt.Println("!")
	cm, _ := clusterManager.NewClusterManager()

	go MasterMetricCollector(cm)

	for {

		host_ctx := "hcp"
		namespace := "hcp"

		host_cfg := cm.Host_config
		//live := cluster.New(host_ctx, host_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})
		live := cluster.New(host_ctx, host_cfg, cluster.Options{})
		fmt.Println("live : ", live)

		ghosts := []*cluster.Cluster{}
		fmt.Println("ghosts: ", ghosts)

		for _, ghost_cluster := range cm.Cluster_list.Items {
			ghost_ctx := ghost_cluster.Name
			fmt.Println("ghost_ctx", ghost_ctx)
			ghost_cfg := cm.Cluster_configs[ghost_ctx]
			fmt.Println("ghost_cfg", ghost_cfg)

			//ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{CacheOptions: cluster.CacheOptions{Namespace: namespace}})
			ghost := cluster.New(ghost_ctx, ghost_cfg, cluster.Options{})
			fmt.Println("ghost", ghost)
			ghosts = append(ghosts, ghost)
			fmt.Println("ghosts", ghosts)
		}

		reshape_cont, _ := reshape.NewController(live, ghosts, namespace, cm)
		fmt.Println("reshpae_cont", reshape_cont)
		// loglevel_cont, _ := logLevel.NewController(live, ghosts, namespace)
		// fmt.Println("loglevel_cont", loglevel_cont)

		m := manager.New()
		m.AddController(reshape_cont)
		// m.AddController(loglevel_cont)

		stop := reshape.SetupSignalHandler()
		fmt.Println("stop", stop)

		if err := m.Start(stop); err != nil {
			log.Fatal(err)
		}
	}

}
func MasterMetricCollector(cm *clusterManager.ClusterManager) {
	// klog.V(4).Info("MasterMetricCollector Called")
	runtime.GOMAXPROCS(runtime.NumCPU())
	INFLUX_IP := os.Getenv("INFLUX_IP")
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
	INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
	INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	// klog.V(5).Info("INFLUX_IP: ", INFLUX_IP)
	// klog.V(5).Info("INFLUX_PORT: ", INFLUX_PORT)
	// klog.V(5).Info("INFLUX_USERNAME: ", INFLUX_USERNAME)
	// klog.V(5).Info("INFLUX_PASSWORD: ", INFLUX_PASSWORD)

	mc := metricCollector.NewMetricCollector(cm, INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)
	// klog.V(2).Info("Created NewMetricCollector Structure")

	// fmt.Println(INFLUX_IP, ":", INFLUX_PORT)
	mc.Influx.CreateDatabase()
	// mc.Influx.CreateMeasurements()

	mc.StartGRPC(GRPC_PORT)

	// mc = &metricCollector.MetricCollector{} //
	// mc.StartGRPC(GRPC_PORT)                 //

}
