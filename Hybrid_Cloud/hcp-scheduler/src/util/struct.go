package util

type TmpEachScore struct {
	Total     int64            // Sum of scores for all clusters
	ScoreList map[string]int64 // key : clusterName, value : SumOfScore
}
