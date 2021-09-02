package DefaultJoin

import (
	"Hybrid_Cluster/hcp-analytic-engine/analyticEngine"
	"Hybrid_Cluster/hcp-apiserver/converter/mappingTable"
	"Hybrid_Cluster/hcp-apiserver/handler"
	"fmt"
)

func DefaultJoin() {
	fmt.Println("[ default Join Process Start ]")
	if policyCheck() {
		fmt.Println("[Option1] Policy exist")
		fmt.Println("--Target Cluster: cluster-1")
		var info mappingTable.ClusterInfo
		info.ClusterName = "cluster-1"
		info.PlatformName = "gke"
		handler.Join(info)
	} else {
		fmt.Println("[Option2] Policy is nonexistent")
		fmt.Println("Call Analytic Engine")
		/*
			url := "http://localhost:8090/HybridctlAnalyticEngine"
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()

			fmt.Println("response Status:", resp.Status)
			fmt.Println("response Headers:", resp.Header)
			if err != nil {
				log.Println(err)
			}
			defer resp.Body.Close()
		*/
		analyticEngine.HybridctlAnalyticEngine()
		fmt.Println("--Target Cluster: cluster-1")
		/* TODO: Analysis result로 대체 */
		var info mappingTable.ClusterInfo
		info.ClusterName = "cluster-1"
		info.PlatformName = "gke"
		handler.Join(info)
	}
}

func policyCheck() bool {
	fmt.Println("-----------------------------------------")
	fmt.Println("Policy Engine Checking")
	fmt.Println("Send Result to User Requirement Checking Module")
	fmt.Println("-----------------------------------------")
	return false
}
