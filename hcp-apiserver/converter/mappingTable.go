package converter

type ClusterInfo struct {
	PlatformName string
	ClusterName  string
}

var AksAPI map[string]string = map[string]string{
	"start": "https://management.azure.com/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/start?api-version=2021-05-01",
	"stop":  "https://management.azure.com/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ContainerService/managedClusters/{resourceName}/stop?api-version=2021-05-01",
}
