package main

import (
	"Hybrid_Cluster/hcp-analytic-engine/analyticEngine"
	"Hybrid_Cluster/hcp-apiserver/converter/mappingTable"
	"Hybrid_Cluster/hcp-apiserver/handler"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/eks"
)

func parser(w http.ResponseWriter, req *http.Request, input interface{}) {
	jsonDataFromHttp, err := ioutil.ReadAll(req.Body)
	fmt.Printf(string(jsonDataFromHttp))
	json.Unmarshal(jsonDataFromHttp, input)
	defer req.Body.Close()
	if err != nil {
		log.Println(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
}

func checkErr(w http.ResponseWriter, err error) {
	if err != nil {
		log.Println(err)
	}
}

// hybridctl join <platformName> <ClusterName>
func join(info mappingTable.ClusterInfo) {

	fmt.Println("---ok---")
	// clusterInfo := mappingTable.ClusterInfo{}
	// var info = mappingTable.ClusterInfo{
	//  PlatformName: clusterInfo["PlatformName"].(string),
	//  ClusterName:  clusterInfo["ClusterName"].(string),
	// }
	// parser(w, req, &clusterInfo)
	// w.Header().Set("Content-Type", "application/json")
	handler.Join(info)
}

func DefaultJoin() {
	fmt.Println("-----------------------------------------")
	fmt.Println("default Join Process Start")
	if policyCheck() {
		fmt.Println("[Option1] Policy exist")
		fmt.Println("--Target Cluster: cluster-1")
		var info mappingTable.ClusterInfo
		info.ClusterName = "cluster-1"
		info.PlatformName = "gke"
		// handler.Join(info)
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
		// handler.Join(info)
	}
}

func policyCheck() bool {
	fmt.Println("-----------------------------------------")
	fmt.Println("Policy Engine Checking")
	fmt.Println("Send Result to User Requirement Checking Module")
	fmt.Println("-----------------------------------------")
	return false
}

func createAddon(w http.ResponseWriter, req *http.Request) {

	var createAddonInput eks.CreateAddonInput

	parser(w, req, &createAddonInput)
	out, err := handler.CreateAddon(createAddonInput)
	// checkErr(w, err)
	var jsonData []byte
	if err != nil {
		log.Println(err)
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Write([]byte(jsonData))
}

func deleteAddon(w http.ResponseWriter, req *http.Request) {

	var deleteAddonInput eks.DeleteAddonInput

	parser(w, req, &deleteAddonInput)
	out, err := handler.DeleteAddon(deleteAddonInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func describeAddon(w http.ResponseWriter, req *http.Request) {

	var describeAddonInput eks.DescribeAddonInput

	parser(w, req, &describeAddonInput)
	out, err := handler.DescribeAddon(describeAddonInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func describeAddonVersions(w http.ResponseWriter, req *http.Request) {

	var describeAddonVersionsInput eks.DescribeAddonVersionsInput

	parser(w, req, &describeAddonVersionsInput)
	out, err := handler.DescribeAddonVersions(describeAddonVersionsInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write(jsonData)

}

func listAddon(w http.ResponseWriter, req *http.Request) {

	var listAddonInput eks.ListAddonsInput

	parser(w, req, &listAddonInput)
	out, err := handler.ListAddon(listAddonInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))

}

func updateAddon(w http.ResponseWriter, req *http.Request) {

	var updateAddonInput eks.UpdateAddonInput

	parser(w, req, &updateAddonInput)
	out, err := handler.UpdateAddon(updateAddonInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func listUpdate(w http.ResponseWriter, req *http.Request) {
	var listUpdateInput eks.ListUpdatesInput

	parser(w, req, &listUpdateInput)
	out, err := handler.ListUpdate(listUpdateInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func describeUpdate(w http.ResponseWriter, req *http.Request) {
	var describeUpdateInput eks.DescribeUpdateInput

	parser(w, req, &describeUpdateInput)
	out, err := handler.DescribeUpdate(describeUpdateInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func listTagsForResource(w http.ResponseWriter, req *http.Request) {
	var listTagsForResourceInput eks.ListTagsForResourceInput

	parser(w, req, &listTagsForResourceInput)
	out, err := handler.ListTagsForResource(listTagsForResourceInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func associateIdentityProviderConfig(w http.ResponseWriter, req *http.Request) {
	var associateIdentityProviderConfigInput eks.AssociateIdentityProviderConfigInput

	parser(w, req, &associateIdentityProviderConfigInput)
	out, err := handler.AssociateIdentityProviderConfig(associateIdentityProviderConfigInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func disassociateIdentityProviderConfig(w http.ResponseWriter, req *http.Request) {
	var disassociateIdentityProviderConfigInput eks.DisassociateIdentityProviderConfigInput

	parser(w, req, &disassociateIdentityProviderConfigInput)
	out, err := handler.DisassociateIdentityProviderConfig(disassociateIdentityProviderConfigInput)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func describeIdentityProviderConfig(w http.ResponseWriter, req *http.Request) {
	var input eks.DescribeIdentityProviderConfigInput

	parser(w, req, &input)
	out, err := handler.DescribeIdentityProviderConfig(input)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func listIdentityProviderConfigs(w http.ResponseWriter, req *http.Request) {
	var input eks.ListIdentityProviderConfigsInput

	parser(w, req, &input)
	out, err := handler.ListIdentityProviderConfigs(input)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func tagResource(w http.ResponseWriter, req *http.Request) {
	var input eks.TagResourceInput

	parser(w, req, &input)
	out, err := handler.TagResource(input)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func untagResource(w http.ResponseWriter, req *http.Request) {
	var input eks.UntagResourceInput

	parser(w, req, &input)
	out, err := handler.UntagResource(input)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func updateClusterConfig(w http.ResponseWriter, req *http.Request) {
	var input eks.UpdateClusterConfigInput

	parser(w, req, &input)
	out, err := handler.UpdateClusterConfig(input)
	checkErr(w, err)
	jsonData, _ := json.Marshal(&out)
	w.Write([]byte(jsonData))
}

func main() {
	// http.HandleFunc("/join", join)
	// // http.HandleFunc("/defaultJoin", defaultJoin)
	http.HandleFunc("/createAddon", createAddon)
	http.HandleFunc("/deleteAddon", deleteAddon)
	http.HandleFunc("/describeAddon", describeAddon)
	http.HandleFunc("/describeAddonVersions", describeAddonVersions)
	http.HandleFunc("/listAddon", listAddon)
	http.HandleFunc("/updateAddon", updateAddon)
	http.HandleFunc("/listUpdate", listUpdate)
	http.HandleFunc("/describeUpdate", describeUpdate)
	http.HandleFunc("/listTagsForResource", listTagsForResource)
	http.HandleFunc("/associateIdentityProvicerConfig", associateIdentityProviderConfig)
	http.HandleFunc("/disassociateIdentityProviderConfig", disassociateIdentityProviderConfig)
	http.HandleFunc("/describeIdentityProviderConfig", describeIdentityProviderConfig)
	http.HandleFunc("/listIdentityProviderConfigs", listIdentityProviderConfigs)
	http.HandleFunc("/tagResource", tagResource)
	http.HandleFunc("/untagResource", untagResource)
	http.HandleFunc("/updateClusterConfig", updateClusterConfig)
	http.ListenAndServe(":8080", nil)
}

/*
*** optionCheck Module ***
- kubernetes Platform Check
*/
// func OptionCheck(w http.ResponseWriter, req *http.Request) {
//  fmt.Println("---Checking Options start---")

// cli := make(map[string]interface{})

// info, err := ioutil.ReadAll(req.Body)
// json.Unmarshal([]byte(jsonDataFromHttp), &cli)
// defer req.Body.Close()
// if err != nil {
//  panic(err)
// }
// handler.JoinHandler(info)

// info := mappingTable.CommandInfo{Cmd: cli["Cmd"].(string), Platform: cli["PlatformName"].(string), ClusterName: cli["ClusterName"].(string)}
// platform 이름 기입여부 체크
// switch info.Platform {
// case "gke", "aks", "eks":
//  printInfo(info.Platform, info.ClusterName)
//  handler.JoinHandler(info)
// default:
//  schedulingPlatform()
//  handler.JoinHandler(info)
// }

// w.Header().Set("Content-Type", "application/json")
// w.WriteHeader(http.StatusOK)
// }

// func printInfo(PlatformName string, clusterName string) {
//  fmt.Println("kubernetes engine Name : ", PlatformName)
//  fmt.Printf("Cluster Name : %s\n", clusterName)
//  fmt.Printf("---Checking Options Done---\n\n")
// }

// func schedulingPlatform() {
//  fmt.Println("---SchedulingPlatform---")
// }
