package util

var AlgorithmList = []string{
	"DRF",
	"Affinity",
}

// var TargetCluster = make(map[string]*fedv1b1.KubeFedCluster)
// var WatchingLevel = policy.GetWatchingLevel()
// var WarningLevel = policy.GetWarningLevel()

var ScoreTable = []*Cluster{
	{
		ClusterInfo: &ClusterInfo{
			ClusterId:   1,
			ClusterName: "cluster1",
		},
		Nodes: []*NodeScore{
			{
				NodeId: 1,
			},
			{
				NodeId: 2,
			},
		},
	}, {
		ClusterInfo: &ClusterInfo{
			ClusterId:   2,
			ClusterName: "cluster2",
		},
		Nodes: []*NodeScore{
			{
				NodeId: 1,
			},
			{
				NodeId: 2,
			},
		},
	}, {
		ClusterInfo: &ClusterInfo{
			ClusterId:   3,
			ClusterName: "cluster3",
		},
		Nodes: []*NodeScore{
			{
				NodeId: 1,
			},
			{
				NodeId: 2,
			},
		},
	}, {
		ClusterInfo: &ClusterInfo{
			ClusterId:   4,
			ClusterName: "cluster4",
		},
		Nodes: []*NodeScore{
			{
				NodeId: 1,
			},
			{
				NodeId: 2,
			},
		},
	},
}
