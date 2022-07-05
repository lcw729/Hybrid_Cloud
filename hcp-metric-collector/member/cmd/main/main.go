package main

import (
	//"context"

	"context"
	"os"

	"github.com/golang/protobuf/ptypes/timestamp"

	"github.com/jinzhu/copier"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"Hybrid_Cloud/hcp-metric-collector/member/pkg/customMetrics"
	"Hybrid_Cloud/hcp-metric-collector/member/pkg/kubeletClient"
	"Hybrid_Cloud/hcp-metric-collector/member/pkg/protobuf"
	"Hybrid_Cloud/hcp-metric-collector/member/pkg/scrap"
	"Hybrid_Cloud/hcp-metric-collector/member/pkg/storage"
	"Hybrid_Cloud/util/clusterManager"

	"fmt"
	"time"
)

func convert(data *storage.Collection, latencyTime string) *protobuf.Collection {
	//klog.V(0).Info("Convert GRPC Data Structure")
	grpc_data := &protobuf.Collection{}

	fmt.Println("-----grpc_data----")
	fmt.Println(grpc_data)
	fmt.Println("----grpcdata end------")
	copier.Copy(grpc_data, data)
	for i := range grpc_data.Metricsbatchs {

		s := int64(data.Metricsbatchs[i].Node.Timestamp.Second())     // from 'int'
		n := int32(data.Metricsbatchs[i].Node.Timestamp.Nanosecond()) // from 'int'

		ts := &timestamp.Timestamp{Seconds: s, Nanos: n}

		mp := &protobuf.MetricsPoint{
			Timestamp:             ts,
			CPUUsageNanoCores:     data.Metricsbatchs[i].Node.CPUUsageNanoCores.String(),
			MemoryUsageBytes:      data.Metricsbatchs[i].Node.MemoryUsageBytes.String(),
			MemoryAvailableBytes:  data.Metricsbatchs[i].Node.MemoryAvailableBytes.String(),
			MemoryWorkingSetBytes: data.Metricsbatchs[i].Node.MemoryWorkingSetBytes.String(),
			NetworkRxBytes:        data.Metricsbatchs[i].Node.NetworkRxBytes.String(),
			NetworkTxBytes:        data.Metricsbatchs[i].Node.NetworkTxBytes.String(),
			FsAvailableBytes:      data.Metricsbatchs[i].Node.FsAvailableBytes.String(),
			FsCapacityBytes:       data.Metricsbatchs[i].Node.FsCapacityBytes.String(),
			FsUsedBytes:           data.Metricsbatchs[i].Node.FsUsedBytes.String(),
			NetworkLatency:        latencyTime,
		}
		grpc_data.Metricsbatchs[i].Node.MP = mp

		//fmt.Println(grpc_data.Metricsbatchs[0].IP)
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.Name)
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.MP.Timestamp.String())
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.MP.Timestamp.Seconds)
		//fmt.Println(grpc_data.Metricsbatchs[0].Node.MP.CpuUsage)
		// fmt.Println(grpc_data.Metricsbatchs[0].Node.MP.MemoryUsage)

		podMetricsPoints := []*protobuf.PodMetricsPoint{}

		for j := range data.Metricsbatchs[i].Pods {
			s := int64(data.Metricsbatchs[i].Pods[j].Timestamp.Second())     // from 'int'
			n := int32(data.Metricsbatchs[i].Pods[j].Timestamp.Nanosecond()) // from 'int'

			ts := &timestamppb.Timestamp{Seconds: s, Nanos: n}

			mp2 := &protobuf.MetricsPoint{
				Timestamp:             ts,
				CPUUsageNanoCores:     data.Metricsbatchs[i].Pods[j].CPUUsageNanoCores.String(),
				MemoryUsageBytes:      data.Metricsbatchs[i].Pods[j].MemoryUsageBytes.String(),
				MemoryAvailableBytes:  data.Metricsbatchs[i].Pods[j].MemoryAvailableBytes.String(),
				MemoryWorkingSetBytes: data.Metricsbatchs[i].Pods[j].MemoryWorkingSetBytes.String(),
				NetworkRxBytes:        data.Metricsbatchs[i].Pods[j].NetworkRxBytes.String(),
				NetworkTxBytes:        data.Metricsbatchs[i].Pods[j].NetworkTxBytes.String(),
				FsAvailableBytes:      data.Metricsbatchs[i].Pods[j].FsAvailableBytes.String(),
				FsCapacityBytes:       data.Metricsbatchs[i].Pods[j].FsCapacityBytes.String(),
				FsUsedBytes:           data.Metricsbatchs[i].Pods[j].FsUsedBytes.String(),
				NetworkLatency:        latencyTime,
			}
			pmp := &protobuf.PodMetricsPoint{
				Name:       data.Metricsbatchs[i].Pods[j].Name,
				Namespace:  data.Metricsbatchs[i].Pods[j].Namespace,
				MP:         mp2,
				Containers: nil,
			}
			podMetricsPoints = append(podMetricsPoints, pmp)

		}
		grpc_data.Metricsbatchs[i].Pods = podMetricsPoints

		//fmt.Println(grpc_data.Metricsbatchs[0].IP)
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].Name)
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].MP.Timestamp.String())
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].MP.Timestamp.Seconds)
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].MP.CpuUsage)
		//fmt.Println(grpc_data.Metricsbatchs[0].Pods[0].MP.MemoryUsage)

	}

	return grpc_data

}
func main() {
	MemberMetricCollector()
}

func MemberMetricCollector() {
	SERVER_IP := os.Getenv("GRPC_SERVER")
	SERVER_PORT := os.Getenv("GRPC_PORT")
	fmt.Println("ClusterMetricCollector Start")
	grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)

	// makedatabase := influxdb.Query{
	// 	Command:  "create database metric",
	// 	Database: "_internal",
	// }

	// grpcClient.Query(makedatabase)

	// var period_int64 int64 = 5
	var period_int64 int64 = 1
	var latencyTime float64 = 0

	for {

		host_config, err := rest.InClusterConfig() // apiserver에 접근해 config파일 읽어오기

		// host_config, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fmt.Println("Host: ", host_config.Host)
		fmt.Println("TLSClientConfig: ", host_config.TLSClientConfig)
		fmt.Println("BearerToken: ", host_config.BearerToken)
		fmt.Println("BearerTokenFile: ", host_config.BearerTokenFile)
		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		host_kubeClient, err := kubernetes.NewForConfig(host_config) // clientSet 생성
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println("host kube client set: ", host_kubeClient)

		Node_list, err := clusterManager.GetNodeList(host_kubeClient) // Node list 읽어오기
		if err != nil {
			fmt.Println(err)
			break
		}

		// fmt.Println("node list: ", Node_list)
		nodes := Node_list.Items // Node list
		fmt.Println("Get Metric Data From Kubelet")
		kubeletClient, _ := kubeletClient.NewKubeletClient()
		fmt.Println("---------------------------print kubeletClient:", kubeletClient)
		data, errs := scrap.Scrap(host_config, kubeletClient, nodes)
		_ = data
		if errs != nil {
			fmt.Println(errs)
			time.Sleep(time.Duration(period_int64) * time.Second)
			continue
		}

		fmt.Println("Convert Metric Data For gRPC")

		latencyTime_string := fmt.Sprintf("%f", latencyTime)

		grpc_data := convert(data, latencyTime_string)

		fmt.Println("[gRPC Start] Send Metric Data")

		// rTime_start := time.Now()

		fmt.Println("--------------------------------------")
		fmt.Println(grpc_data)
		fmt.Println("--------------------------------------")

		r, err := grpcClient.SendMetrics(context.TODO(), grpc_data)
		// rTime_end := time.Since(rTime_start)

		// latencyTime = rTime_end.Seconds() - r.ProcessingTime

		// fmt.Println("rTime_end: ", rTime_end)
		// fmt.Println("rProcessingTime: ", r.ProcessingTime)

		if err != nil {
			//fmt.Println("check")
			fmt.Println("could not connect : ", err)
			time.Sleep(time.Duration(period_int64) * time.Second)
			//fmt.Println("check2")
			continue
		}
		fmt.Println("[gRPC End] Send Metric Data")

		// period_int64 := r.Tick
		// _ = data

		fmt.Println("[http Start] Post Metric Data to Custom Metric Server")
		token := host_config.BearerToken
		host := host_config.Host
		client := host_kubeClient
		fmt.Println("host: ", host)
		fmt.Println("token: ", token)
		fmt.Println("client: ", client)

		customMetrics.AddToPodCustomMetricServer(data, token, host)
		customMetrics.AddToDeployCustomMetricServer(data, token, host, client)
		fmt.Println("[http End] Post Metric Data to Custom Metric Server")

		period_int64 = r.Tick

		if period_int64 > 0 && err == nil {

			//fmt.Println("period : ",time.Duration(period_int64))
			fmt.Println("Wait ", time.Duration(period_int64)*time.Second, "...")
			time.Sleep(time.Duration(period_int64) * time.Second)
		} else {
			fmt.Println("--- Fail to get period")
			time.Sleep(5 * time.Second)
		}
	}
}
