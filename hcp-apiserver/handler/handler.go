package handler

import (
	"Hybrid_Cluster/hybridctl/util"
	"fmt"
	"net/http"
)

func AksStart(util.EksAPIParameter) (*http.Response, error) {
	httpPostUrl := "https://management.azure.com/subscriptions/ccfc0c6c-d3c6-4de2-9a6c-c09ca498ff73/resourceGroups/hcpResourceGroup/providers/Microsoft.ContainerService/managedClusters/aks-master/start?api-version=2021-05-01"

	request, _ := http.NewRequest("POST", httpPostUrl, nil)
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	return response, err
}
