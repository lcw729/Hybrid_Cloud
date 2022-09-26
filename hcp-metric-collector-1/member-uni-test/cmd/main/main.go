package main

import (
	//"context"

	"context"
	"io/ioutil"
	"net"
	"os"

	"github.com/golang/protobuf/ptypes/timestamp"

	// flowcontrol "k8s.io/client-go/util/flowcontrol"

	"github.com/jinzhu/copier"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/klog"

	"Hybrid_Cloud/hcp-metric-collector/member-uni-test/pkg/customMetrics"
	"Hybrid_Cloud/hcp-metric-collector/member/pkg/kubeletClient"
	"Hybrid_Cloud/hcp-metric-collector/member/pkg/protobuf"
	"Hybrid_Cloud/hcp-metric-collector/member/pkg/scrap"
	"Hybrid_Cloud/hcp-metric-collector/member/pkg/storage"
	"Hybrid_Cloud/util/clusterManager"

	"fmt"
	"time"
)

type NodeList struct {
	ApiVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Items      []Node
}

type Node struct {
	Metadata Metadata   `json:"metadata"`
	Status   NodeStatus `json:"status"`
}

type ResourceList map[string]string
type NodeStatus struct {
	Capacity    ResourceList `json:"capacity"`
	Allocatable ResourceList `json:"allocatable"`
}

type Metadata struct {
	Name            string            `json:"name"`
	GenerateName    string            `json:"generateName"`
	ResourceVersion string            `json:"resourceVersion"`
	Labels          map[string]string `json:"labels"`
	Annotations     map[string]string `json:"annotations"`
	Uid             string            `json:"uid"`
}

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

func InClusterConfig() (*rest.Config, error) {
	const (
		tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")

	fmt.Println("host: ", host)
	fmt.Println("port: ", port)
	fmt.Println("host len: ", len(host), " port len: ", len(port))

	if len(host) == 0 || len(port) == 0 {
		fmt.Println("host: ", host)
		fmt.Println("port: ", port)
		fmt.Println("host len: ", len(host), " port len: ", len(port))
	}

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		fmt.Println(err)

	}

	tlsClientConfig := rest.TLSClientConfig{}

	if _, err := certutil.NewPool(rootCAFile); err != nil {
		klog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
		// os.Exit(3) // uni log 확인하려면 아예 꺼져야 함.
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	tlsClientConfig.CAFile = rootCAFile

	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,
	}, nil
}

func main() {
	os.Exit(1)
	MemberMetricCollector()
}

// func getNodes() (*NodeList, error) {
// 	var nodeList NodeList

// 	request := &http.Request{
// 		Header: make(http.Header),
// 		Method: http.MethodGet,
// 		URL: &url.URL{
// 			Host:   apiHost,
// 			Path:   nodesEndpoint,
// 			Scheme: "http",
// 		},
// 	}
// 	request.Header.Set("Accept", "application/json, */*")

// 	resp, err := http.DefaultClient.Do(request)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = json.NewDecoder(resp.Body).Decode(&nodeList)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &nodeList, nil
// }

func createFile(path string) {
	// check if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if isError(err) {
			return
		}
		defer file.Close()
	}

	fmt.Println("File Created Successfully", path)
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}
func writeFile(path string) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// Write some text line-by-line to file.
	_, err = file.WriteString("Hello \n")
	if isError(err) {
		return
	}
	_, err = file.WriteString("World \n")
	if isError(err) {
		return
	}

	// Save file changes.
	err = file.Sync()
	if isError(err) {
		return
	}

	fmt.Println("File Updated Successfully.")
}

func MemberMetricCollector() {
	// SERVER_IP := os.Getenv("GRPC_SERVER")
	// SERVER_IP := "115.94.141.62"
	SERVER_IP := "10.0.5.83"
	// SERVER_PORT := os.Getenv("GRPC_PORT")
	SERVER_PORT := "32051"
	// SERVER_PORT := "2051"
	fmt.Println("ClusterMetricCollector Start")
	grpcClient := protobuf.NewGrpcClient(SERVER_IP, SERVER_PORT)
	fmt.Println(grpcClient)
	createFile("")

	// makedatabase := influxdb.Query{
	// 	Command:  "create database metric",
	// 	Database: "_internal",
	// }

	// grpcClient.Query(makedatabase)

	// var period_int64 int64 = 5
	var period_int64 int64 = 1
	var latencyTime float64 = 0

	for {
		//여기부터
		// host_config, err := uni_Config() // apiserver에 접근해 config파일 읽어오기
		host_config, _ := InClusterConfig() // apiserver에 접근해 config파일 읽어오기

		// fmt.Println(host_config)

		host_kubeClient, err := kubernetes.NewForConfig(host_config) // clientSet 생성
		if err != nil {
			fmt.Println(err)

		}
		//여기까지

		// fmt.Println("create client set: ", host_kubeClient)
		//이 노드 리스트 이부분에서 에러 발생 .04.04
		Node_list, err := clusterManager.GetNodeList(host_kubeClient) // Node list 읽어오기
		if err != nil {
			fmt.Println(err)

		}

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

		fmt.Print(latencyTime_string)

		grpc_data := convert(data, latencyTime_string)

		fmt.Println("[gRPC Start] Send Metric Data")

		// // rTime_start := time.Now()

		fmt.Println("-------------------!-------------------")
		fmt.Println(host_config)
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
			// break
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
