package handler

import (
	"Hybrid_Cluster/hcp-apiserver/converter"
	auth "Hybrid_Cluster/hcp-apiserver/util"
	"Hybrid_Cluster/hybridctl/util"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func AksStart(input util.EksAPIParameter) (*http.Response, error) {
	api := converter.AksAPI["start"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ResourceName)
	fmt.Println(api)
	hosturl := api
	response, err := auth.AuthorizationAndHTTP("POST", hosturl)
	return response, err
}

func AksStop(input util.EksAPIParameter) (*http.Response, error) {
	api := converter.AksAPI["stop"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ResourceName)
	hosturl := api
	response, err := auth.AuthorizationAndHTTP("POST", hosturl)
	return response, err
}

func AksRotateCerts(input util.EksAPIParameter) (*http.Response, error) {
	api := converter.AksAPI["rotateCerts"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ResourceName)
	fmt.Println(api)
	hosturl := api
	response, err := auth.AuthorizationAndHTTP("POST", hosturl)
	return response, err
}

func AksGetOSoptions(input util.EksAPIParameter) (*http.Response, error) {
	api := converter.AksAPI["getOSoptions"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{location}", input.Location)
	hosturl := api
	fmt.Println(api)
	response, err := auth.AuthorizationAndHTTP("GET", hosturl)
	return response, err
}
