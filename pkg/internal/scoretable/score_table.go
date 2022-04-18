package scoretable

import (
	"math"
	"sort"
)

const (
	// MaxNodeScore is the maximum score a Score plugin is expected to return.
	MaxNodeScore int32 = 100

	// MinNodeScore is the minimum score a Score plugin is expected to return.
	MinNodeScore int32 = 0

	// MaxTotalScore is the maximum total score.
	MaxTotalScore int32 = math.MaxInt32
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

func NewScoreTable(size int) *map[string]float32 {
	score_table := make(map[string]float32, size)

	return &score_table
}

func SortScore(score_table map[string]float32) ScoreTable {
	sorted := make(ScoreTable, len(score_table))

	for cluster, score := range score_table {
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
