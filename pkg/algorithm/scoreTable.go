package algorithm

import "sort"

type Score struct {
	Key   string
	Value float32
}

type ScoreTable []Score

func (s ScoreTable) Len() int {
	return len(s)
}

func (s ScoreTable) Less(i, j int) bool {
	return s[i].Value < s[j].Value
}

func (s ScoreTable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var score_table = make(map[string]float32)

func SortScore() ScoreTable {
	sorted := make(ScoreTable, len(score_table))
	var i int
	for key, value := range score_table {
		sorted[i] = Score{key, value}
		i++
	}
	sort.Sort(sorted)

	return sorted
}
