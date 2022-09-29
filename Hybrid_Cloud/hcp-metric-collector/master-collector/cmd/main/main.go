package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/KETI-Hybrid/hcp-metric-collector-m-v1/pkg/metricCollector"

	"github.com/KETI-Hybrid/hcp-pkg/util/clusterManager"

	"k8s.io/klog"

	"github.com/jinzhu/copier"
	"k8s.io/client-go/rest"
	"k8s.io/sample-controller/pkg/signals"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
	genericclient "sigs.k8s.io/kubefed/pkg/client/generic"
)

const (
	GRPC_PORT = "2051"
)

var prev_length int = 0
var mc *metricCollector.MetricCollector

func main() {
	var wg sync.WaitGroup

	cm, _ := clusterManager.NewClusterManager()

	stopCh := signals.SetupSignalHandler()

	wg.Add(2)
	go MasterMetricCollector(cm, stopCh)
	go reshapeCluster(stopCh)
	wg.Wait()
}

func reshapeCluster(stopCh <-chan struct{}) {
	for {
		host_config, err := rest.InClusterConfig()
		if err != nil {
			<-stopCh
		}

		namespace := "kube-federation-system"
		host_client := genericclient.NewForConfigOrDie(host_config)
		tempClusterList := &fedv1b1.KubeFedClusterList{}

		err = host_client.List(context.TODO(), tempClusterList, namespace, &client.ListOptions{})

		if err != nil {
			fmt.Printf("Error retrieving list of federated clusters: %+v\n", err)
			<-stopCh
		}
		temp_length := len(tempClusterList.Items)

		if temp_length != prev_length {
			klog.Infof("temp_length : %d, prev_length", temp_length, prev_length)
			newCm, err := clusterManager.NewClusterManager()

			if err != nil {
				klog.Errorln(err)
				<-stopCh
			}

			copier.Copy(mc.ClusterManager, newCm)
			prev_length = temp_length
		}
	}
}

func MasterMetricCollector(cm *clusterManager.ClusterManager, stopCh <-chan struct{}) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	INFLUX_IP := os.Getenv("INFLUX_IP")
	INFLUX_PORT := os.Getenv("INFLUX_PORT")
	INFLUX_USERNAME := os.Getenv("INFLUX_USERNAME")
	INFLUX_PASSWORD := os.Getenv("INFLUX_PASSWORD")

	mc = metricCollector.NewMetricCollector(cm, INFLUX_IP, INFLUX_PORT, INFLUX_USERNAME, INFLUX_PASSWORD)
	mc.Influx.CreateDatabase()
	mc.StartGRPC(GRPC_PORT)
}
