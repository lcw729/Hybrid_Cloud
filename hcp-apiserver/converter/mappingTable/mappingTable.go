package mappingTable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ClusterInfo struct {
	PlatformName string
	ClusterName  string
}

// type CommandInfo struct {
// 	Cmd         string
// 	Platform    string
// 	ClusterName string
// }

type aksData struct {
	ResourceGroupName string
	ResourceName      string
	agentPoolName     string
	ApiVersion        string
}

var aks aksData
var apiMap map[ClusterInfo]string = make(map[ClusterInfo]string)

func GetInfo(info ClusterInfo) string {
	return apiMap[info]
}

// 해당 명령어 API 반환
func AksGetCredential(info ClusterInfo) string {
	fmt.Println("---AKS GetCredential API---")
	httpGetUrl := "http://10.0.5.43:8080/aksGetCredential"
	response, err := http.Get(httpGetUrl)

	if err != nil {
		fmt.Println("Data request failed")
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal([]byte(body), &aks)
	if info.ClusterName == aks.ResourceName {
		apiMap[ClusterInfo{info.PlatformName, info.ClusterName}] = "resourceGroups/" + aks.ResourceGroupName + "/providers/Microsoft.ContainerService/managedClusters/" + aks.ResourceName + "/listClusterUserCredential?api-version=" + aks.ApiVersion
		api := apiMap[ClusterInfo{info.PlatformName, info.ClusterName}]
		fmt.Println("---End AKS GetCredential API---")
		return api
	} else {
		fmt.Printf("-> NotFound clusterName : %s\n", info.ClusterName)
		return ""
	}
}
