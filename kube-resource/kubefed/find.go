package kubefed

import "Hybrid_Cloud/util/clusterManager"

func Iskubefedcluster(clustername string) bool {
	cm, _ := clusterManager.NewClusterManager()
	list := cm.Cluster_list
	for _, i := range list.Items {
		if i.Name == clustername {
			return true
		}
	}
	return false
}
