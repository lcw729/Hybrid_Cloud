package scoretable

import (
	"sort"
)

type Score struct {
	Cluster string
	Score   float32
}

type ScoreTable []Score

func (s ScoreTable) Len() int {
	return len(s)
}

func (s ScoreTable) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func (s ScoreTable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type TargetClustersScoreTable map[string]float32 // 정렬하고 싶은 map

func NewScoreTable(clusterList *[]string) TargetClustersScoreTable {
	score_table := make(map[string]float32, len(*clusterList))
	for _, i := range *clusterList {
		score_table[i] = 0
	}
	return score_table
}

func (s *TargetClustersScoreTable) SortScore() ScoreTable {
	sorted := make(ScoreTable, len(*s))

	for cluster, score := range *s {
		sorted = append(sorted, Score{cluster, score})
	}
	sort.Sort(sort.Reverse(sorted))

	return sorted
}

/*
func (s *ClusterScoreTable) SortCluster() []string {
	sorted := make([]string, 0, len(*s))

	for cluster := range *s {
		sorted = append(sorted, cluster)
	}
	sort.Strings(sorted)

	return sorted
}
*/
