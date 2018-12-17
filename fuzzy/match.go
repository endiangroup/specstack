package fuzzy

import (
	"math"
	"sort"

	"github.com/antzucaro/matchr"
)

const (
	DistanceThreshold = 0.75
	MinThreshold      = 0.25
	AdjacentThreshold = 0.05
)

type Match struct {
	Term string
	Rank float64
}

func (r Match) Negligible() bool {
	return r.Rank < MinThreshold
}

func Rank(comparison string, indexes []string) []Match {
	output := make([]Match, len(indexes))
	for i, index := range indexes {
		output[i] = Match{
			Term: index,
			Rank: Strcmp(comparison, index),
		}
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].Rank > output[j].Rank
	})
	return output
}

func Strcmp(a, b string) float64 {
	maxLen := math.Max(float64(len(a)), float64(len(b)))
	levDist := matchr.DamerauLevenshtein(a, b)
	return (1 - (float64(levDist) / maxLen))
}

func Equivalent(a, b Match) bool {
	return math.Abs(a.Rank-b.Rank) < AdjacentThreshold
}
