package handler

import (
	auth "Hybrid_Cluster/hcp-apiserver/util"
	"Hybrid_Cluster/hybridctl/util"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

//var azureAuth = auth.GetAzureAuth()

func AksStart(input util.EksAPIParameter) (*http.Response, error) {
	hosturl := "https://management.azure.com/subscriptions/" + os.Getenv("SubscriptionId") + "/resourceGroups/" + input.ResourceGroupName + "/providers/Microsoft.ContainerService/managedClusters/" + input.ResourceName + "/start?api-version=2021-05-01"
	response, err := AuthorizationAndPost(hosturl)
	return response, err
}

func AksStop(input util.EksAPIParameter) (*http.Response, error) {
	hosturl := "https://management.azure.com/subscriptions/" + os.Getenv("SubscriptionId") + "/resourceGroups/" + input.ResourceGroupName + "/providers/Microsoft.ContainerService/managedClusters/" + input.ResourceName + "/stop?api-version=2021-05-01"
	response, err := AuthorizationAndPost(hosturl)
	return response, err
}

func AuthorizationAndPost(hosturl string) (*http.Response, error) {
	params := url.Values{}
	params.Add("resource", `https://management.azure.com/`)
	body := strings.NewReader(params.Encode())
	request, _ := http.NewRequest("POST", hosturl, body)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Add("Authorization", "Bearer "+auth.GetBearer().Access_token)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("err2 : ", err)
	} else {
		fmt.Println(response.Status)
	}
	return response, err
}
