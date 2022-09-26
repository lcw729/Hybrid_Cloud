package main

import (
	//"context"

	"context"
	"net"
	"os"

	"github.com/golang/protobuf/ptypes/timestamp"

	// flowcontrol "k8s.io/client-go/util/flowcontrol"

	"github.com/jinzhu/copier"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"Hybrid_Cloud/hcp-metric-collector/member-uni-test/pkg/customMetrics"
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

func InClusterConfig() (*rest.Config, error) {
	const (
		tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host := "10.96.0.1"

	port := "443"

	if len(host) == 0 || len(port) == 0 {
		os.Exit(3)
		return nil, rest.ErrNotInCluster
	}

	// token, err := ioutil.ReadFile(tokenFile)
	// if err != nil {
	// 	os.Exit(3)
	// }
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6InN6c3RGdmJrOW9zN09ldGJWUkhndDUxUnVPc1RkcmZqMVREVTFOZVFtVnMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJoY3AiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiaGNwLW1ldHJpYy1jb2xsZWN0b3ItdG9rZW4tNm1nYnMiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiaGNwLW1ldHJpYy1jb2xsZWN0b3IiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIwN2ZkY2M0Zi0zYzVjLTQ4NTYtYWRkNS05ZWNhNzA5YjA2ZGEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6aGNwOmhjcC1tZXRyaWMtY29sbGVjdG9yIn0.akqrcvzYXCEW5ixBuIw79yfLcYYtMnsOa0UhZa9LPzvHYAoNYFY28m-zKbKVyN3XaavyTkgfqaeKq-AohEbKzfZH3_7olRoU_6E4ltASrs9IHppXKNxgWDzkU25xyZa9Kskr-hcIjefyZmepSj1kNrJfcuL4fl7OWo33AIrPfWF6SRO7fSX2_nUcMWjgr1SkOQzLSbbYGU6oT5EGRY8x3nGyeA4jrjMp_-rJJliGGi5AcNVwskmBwlP4qBFpfDWSApCecJ_YOP5f5DXYn3dpfcxDdQiMi4Af4vLDeQYCMOg77slWZz9SIM90pDCLw3Tz4i4hu3SZqxpiSN5z9N2naQ"

	tlsClientConfig := rest.TLSClientConfig{}

	// if _, err := certutil.NewPool(rootCAFile); err != nil {
	// 	klog.Errorf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	// 	os.Exit(3) // uni log 확인하려면 아예 꺼져야 함.
	// } else {
	// 	tlsClientConfig.CAFile = rootCAFile
	// }

	tlsClientConfig.CAFile = rootCAFile

	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,
	}, nil
}

// func uni_Config() (*rest.Config, error) {

// 	return &rest.Config{
// 		Host: "https://10.96.0.1:443",
// 		// TLSClientConfig: rest.sanitizedTLSClientConfig{Insecure: false, ServerName: "", CertFile: "", KeyFile: "", CAFile: "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt", CertData: []uint8(nil), KeyData: []uint8(nil), CAData: []uint8(nil), NextProtos: []string(nil)},
// 		// BearerToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6InN6c3RGdmJrOW9zN09ldGJWUkhndDUxUnVPc1RkcmZqMVREVTFOZVFtVnMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJkZWZhdWx0Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZWNyZXQubmFtZSI6ImRlZmF1bHQtdG9rZW4tcHNkaHEiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiZGVmYXVsdCIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VydmljZS1hY2NvdW50LnVpZCI6IjljMzg1NWFkLTY2MTEtNDkyNi04YjU5LThiNzk0OTVjMWI5MyIsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmRlZmF1bHQifQ.SY2sLOgRlxUeYdpErSl5QtjFCfM2ScuxpjoaBhxS30Cs7DQksmbeSrSwPUPEw3GJWm2wdgqnniUh7kdLcBnwtfU0sY2XYXJo_gWmnhI2ZBFlwNWDdooQNDH2WDPtbgHfyLWd2e1tzD0MkXlLdGYLt9hRbExTu3VTO-zDCfeWM2d01Yw7v2Ipu2XX8yFcc_yMuleNXnD0Va8YQ_enqjyDrbSTY5aqwtS_ps3R-CdTJduO4gld-LqDdEtIV4jSg7tfAianwGwTB0-QhCrjbU5tBw-b5jQ2FDsitmAO4_RUJfFZVVJj-96LVgDPlhIqSNM99BVL7pgE2Ew05h2g5LMBJw",
// 		BearerToken:     "eyJhbGciOiJSUzI1NiIsImtpZCI6InN6c3RGdmJrOW9zN09ldGJWUkhndDUxUnVPc1RkcmZqMVREVTFOZVFtVnMifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJoY3AiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlY3JldC5uYW1lIjoiaGNwLW1ldHJpYy1jb2xsZWN0b3ItdG9rZW4tNm1nYnMiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC5uYW1lIjoiaGNwLW1ldHJpYy1jb2xsZWN0b3IiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiIwN2ZkY2M0Zi0zYzVjLTQ4NTYtYWRkNS05ZWNhNzA5YjA2ZGEiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6aGNwOmhjcC1tZXRyaWMtY29sbGVjdG9yIn0.akqrcvzYXCEW5ixBuIw79yfLcYYtMnsOa0UhZa9LPzvHYAoNYFY28m-zKbKVyN3XaavyTkgfqaeKq-AohEbKzfZH3_7olRoU_6E4ltASrs9IHppXKNxgWDzkU25xyZa9Kskr-hcIjefyZmepSj1kNrJfcuL4fl7OWo33AIrPfWF6SRO7fSX2_nUcMWjgr1SkOQzLSbbYGU6oT5EGRY8x3nGyeA4jrjMp_-rJJliGGi5AcNVwskmBwlP4qBFpfDWSApCecJ_YOP5f5DXYn3dpfcxDdQiMi4Af4vLDeQYCMOg77slWZz9SIM90pDCLw3Tz4i4hu3SZqxpiSN5z9N2naQ",
// 		BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
// 	}, nil
// }

func main() {

	MemberMetricCollector()
}

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

func writeFile(path string, stringerr string) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// Write some text line-by-line to file.
	_, err = file.WriteString(stringerr)
	if isError(err) {
		return
	}
	_, err = file.WriteString(stringerr)
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
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
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

	// makedatabase := influxdb.Query{
	// 	Command:  "create database metric",
	// 	Database: "_internal",
	// }

	// grpcClient.Query(makedatabase)

	// var period_int64 int64 = 5
	var period_int64 int64 = 1
	var latencyTime float64 = 0
	createFile("InClusterConfig.txt")
	createFile("NewForConfig.txt")
	createFile("Nodelist.txt")

	for {

		// host_config, err := uni_Config() // apiserver에 접근해 config파일 읽어오기
		host_config, err := rest.InClusterConfig() // apiserver에 접근해 config파일 읽어오기
		if err != nil {
			writeFile("InClusterConfig.txt", err.Error())
		}

		// fmt.Println("config file: ", host_config)

		for i := 0; i < 1000; i++ {
			fmt.Println(i, " times test")
			time.Sleep(time.Second * 3)
		}

		host_kubeClient, err := kubernetes.NewForConfig(host_config) // clientSet 생성
		if err != nil {
			writeFile("NewForConfig.txt", err.Error())
		}

		for i := 0; i < 1000; i++ {
			fmt.Println(i, " times test")
			time.Sleep(time.Second * 3)
		}

		// fmt.Println("create client set: ", host_kubeClient)
		//이 노드 리스트 이부분에서 에러 발생 .04.04
		Node_list, err := clusterManager.GetNodeList(host_kubeClient) // Node list 읽어오기
		if err != nil {
			writeFile("Nodelist.txt", err.Error())
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
