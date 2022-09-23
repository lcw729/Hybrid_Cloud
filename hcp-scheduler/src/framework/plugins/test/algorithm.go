package test

import (
	policy "github.com/KETI-Hybrid/hcp-pkg/hcp-resource/hcppolicy"
	clusterManager "github.com/KETI-Hybrid/hcp-pkg/util/clusterManager"

	"k8s.io/klog"
	fedv1b1 "sigs.k8s.io/kubefed/pkg/apis/core/v1beta1"
)

// var AlgorithmMap = map[string]func() bool{
// 	"DRF":      DRF,
// 	"Affinity": Affinity,
// }

var TargetCluster = make(map[string]*fedv1b1.KubeFedCluster)
var cm, _ = clusterManager.NewClusterManager()

func WatchingLevelCalculator() {
	klog.Infoln("-----------------------------------------")
	klog.Infoln("[step 2] Get Policy - watching level & warning level")
	watching_level := policy.GetWatchingLevel(*cm.HCPPolicy_Client)
	klog.Infoln("< Watching Level > \n", watching_level)
	// 각 클러스터의 watching level 계산하고 warning level 초과 시 targetCluster에 추가
	warning_level := policy.GetWarningLevel(*cm.HCPPolicy_Client)
	klog.Infoln("< Warning  Level > \n", warning_level)
	klog.Infoln("-----------------------------------------")
	klog.Infoln("[step 3] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	klog.Infoln("[step 4] Calculate watching level")

	cm, err := clusterManager.NewClusterManager()
	if err != nil {
		klog.Error(err)
		return
	}
	for _, cluster := range cm.Cluster_list.Items {
		klog.Infoln(cluster.Name)
		TargetCluster[cluster.Name] = &cluster
	}
	klog.Infoln(TargetCluster)
	// cluster := &util.ClusterInfo{
	// 	ClusterId:   1,
	// 	ClusterName: "cluster1",
	// }
	// if !appendTargetCluster(cluster) {
	// 	klog.Infof("%d exist already\n", cluster.ClusterId)
	// } else {
	// 	klog.Infoln("ok")
	// }

	// cluster = &util.ClusterInfo{
	// 	ClusterId:   2,
	// 	ClusterName: "cluster2",
	// }
	// if !appendTargetCluster(cluster) {
	// 	klog.Infof("%d exist already\n", cluster.ClusterId)
	// } else {
	// 	klog.Infoln("ok")
	// }
}

// func appendTargetCluster(cluster *util.ClusterInfo) bool {
// 	var check bool = false
// 	for _, c := range util.TargetCluster {
// 		if c.ClusterId == cluster.ClusterId {
// 			check = true
// 			break
// 		}
// 	}
// 	if !check {
// 		util.TargetCluster = append(util.TargetCluster, cluster)
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func ResourceConfigurationManagement() {
// 	// targetCluster := WatchingLevelCalculator()
// 	WatchingLevelCalculator()
// 	// klog.Infoln("[step 5] Start ResourceConfiguration")
// 	// for index, cluster := range targetCluster {
// 	// 	klog.Infoln("Index : ", index, "\nClusterId : ", cluster.ClusterId, "\nClusterName : ", cluster.ClusterName)
// 	// }
// }

// // 최적 배치 알고리즘
// func Affinity(clusterList *[]string) string {
// 	klog.Infoln("---------------------------------------------------------------")
// 	klog.Infoln("Affinity Calculator Called")
// 	klog.Infoln("[step 2] Get MultiMetric")
// 	// monitoringEngine.MetricCollector()
// 	klog.Infoln("[step 3-1] Start analysis Resource Affinity")
// 	// score_table := scoretable.NewScoreTable(clusterList)
// 	// for _, i := range *clusterList {
// 	// 	score_table[i] = rand.Float32()
// 	// }
// 	// result := score_table.SortScore()
// 	klog.Infoln("[step 3-2] Send analysis result to Scheduler [Target Cluster]")
// 	klog.Infoln("---------------------------------------------------------------")
// 	// klog.Infoln(score_table)
// 	// return result[0].Cluster
// }

func DRF() bool {
	klog.Infoln("DRF Math operation Called")
	klog.Infoln("-----------------------------------------")
	klog.Infoln("[step 2] Get MultiMetric")
	// monitoringEngine.MetricCollector()
	klog.Infoln("-----------------------------------------")
	klog.Infoln("[step 3-1] Handling Math Operation")
	klog.Infoln("[step 3-2] Search Pod Fit Resources")
	klog.Infoln("[step 3-3] Schedule Decision")
	klog.Infoln("---------------------------------------------------------------")
	return true
}

/*
// 배치 알고리즘에 설정 값에 따라 알고리즘 변경
func OptimalArrangementAlgorithm() map[string]float32 {
	klog.Infoln("[step 1] Get Policy - algorithm")
	algo := policy.GetAlgorithm()
	if algo != "" {
		klog.Infoln(algo)
		switch algo {
		case "DRF":
			return DRF()
		case "Affinity":
			return Affinity()
		default:
			klog.Infoln("there is no such algorithm.")
			return false
		}
	} else {
		klog.Infoln("there is no such algorithm.")
		return false
	}
}
*/

/*
// 가장 점수가 높은 Cluster, Node 확인
func OptimalNodeSelector() (*util.Cluster, *util.NodeScore) {
	max := 0
	cluster := util.ScoreTable[0]
	node := util.ScoreTable[0].Nodes[0]
	for _, c := range util.ScoreTable {
		for _, n := range c.Nodes {
			if int(n.Score) > max {
				max = int(n.Score)
				cluster = c
				node = n
			}
		}
	}
	return cluster, node
}
*/
