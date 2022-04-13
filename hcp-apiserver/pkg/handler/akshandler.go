package handler

import (
	"Hybrid_Cloud/hcp-apiserver/pkg/converter"
	util "Hybrid_Cloud/hcp-apiserver/pkg/util"

	"fmt"
	"net/http"
	"os"
	"strings"
)

func AKSStart(input util.AKSAPIParameter) (*http.Response, error) {
	api := converter.AKSAPI["start"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	fmt.Println(os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ClusterName)
	fmt.Println(api)
	hosturl := api
	response, err := util.AuthorizationAndHTTP("POST", hosturl, nil)
	return response, err
}

func AKSStop(input util.AKSAPIParameter) (*http.Response, error) {
	api := converter.AKSAPI["stop"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ClusterName)
	hosturl := api
	response, err := util.AuthorizationAndHTTP("POST", hosturl, nil)
	return response, err
}

func AKSRotateCerts(input util.AKSAPIParameter) (*http.Response, error) {
	api := converter.AKSAPI["rotateCerts"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ClusterName)
	fmt.Println(api)
	hosturl := api
	response, err := util.AuthorizationAndHTTP("POST", hosturl, nil)
	return response, err
}

func AKSGetOSoptions(input util.AKSAPIParameter) (*http.Response, error) {
	api := converter.AKSAPI["getOSoptions"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{location}", input.Location)
	hosturl := api
	fmt.Println(api)
	response, err := util.AuthorizationAndHTTP("GET", hosturl, nil)
	return response, err
}

func MaintenanceconfigurationCreateOrUpdate(input util.AKSAPIParameter) (*http.Response, error) {
	api := converter.AKSAPI["maintenanceconfigurationCreate/Update"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ClusterName)
	api = strings.ReplaceAll(api, "{configName}", input.ConfigName)
	fmt.Println(api)
	hosturl := api
	response, err := util.AuthorizationAndHTTP("PUT", hosturl, input.ConfigFile)
	return response, err
}

func MaintenanceconfigurationDelete(input util.AKSAPIParameter) (*http.Response, error) {
	api := converter.AKSAPI["maintenanceconfigurationDelete"]
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ClusterName)
	api = strings.ReplaceAll(api, "{configName}", input.ConfigName)
	hosturl := api
	fmt.Println(api)
	response, err := util.AuthorizationAndHTTP("DELETE", hosturl, nil)
	return response, err
}

func MaintenanceconfigurationList(input util.AKSAPIParameter) (*http.Response, error) {
	api := converter.AKSAPI["maintenanceconfigurationList"]
	fmt.Println(input)
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ClusterName)
	fmt.Println(api)
	hosturl := api
	response, err := util.AuthorizationAndHTTP("GET", hosturl, nil)
	return response, err
}

func MaintenanceconfigurationShow(input util.AKSAPIParameter) (*http.Response, error) {
	api := converter.AKSAPI["maintenanceconfigurationShow"]
	fmt.Println(input)
	api = strings.ReplaceAll(api, "{subscriptionId}", os.Getenv("SubscriptionId"))
	api = strings.ReplaceAll(api, "{resourceGroupName}", input.ResourceGroupName)
	api = strings.ReplaceAll(api, "{resourceName}", input.ClusterName)
	api = strings.ReplaceAll(api, "{configName}", input.ConfigName)
	hosturl := api
	response, err := util.AuthorizationAndHTTP("GET", hosturl, nil)
	return response, err
}
