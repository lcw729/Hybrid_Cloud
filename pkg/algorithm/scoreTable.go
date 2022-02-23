package algorithm

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

var score_table = make(map[string]float32) // 정렬하고 싶은 map

func SortScore() ScoreTable {
	sorted := make(ScoreTable, len(score_table))

	for cluster, score := range score_table {
		sorted = append(sorted, Score{cluster, score})
	}
	sort.Sort(sort.Reverse(sorted))

	return sorted
}

func SortCluster() []string {
	sorted := make([]string, 0, len(score_table))

	for cluster := range score_table {
		sorted = append(sorted, cluster)
	}
	sort.Strings(sorted)

	return sorted
}
