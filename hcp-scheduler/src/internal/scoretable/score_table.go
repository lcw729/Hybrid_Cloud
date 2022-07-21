package scoretable

import (
	"math"
	"sort"
)

const (
	// MaxNodeScore is the maximum score a Score plugin is expected to return.
	MaxNodeScore int64 = 100

	// MinNodeScore is the minimum score a Score plugin is expected to return.
	MinNodeScore int64 = 0

	// MaxTotalScore is the maximum total score.
	MaxTotalScore int64 = math.MaxInt64

	MaxCount
)

type ClusterScoreList []ClusterScore

type ClusterScore struct {
	Cluster       string
	NodeScoreList NodeScoreList
	Score         float32
}

type NodeScoreList []NodeScore

type NodeScore struct {
	Name  string
	Score int64
}

type ScoreTable ClusterScoreList

func (s ScoreTable) Len() int {
	return len(s)
}

func (s ScoreTable) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

func (s ScoreTable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *ScoreTable) SortScore() {
	sort.Sort(sort.Reverse(*s))
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
