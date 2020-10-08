package fuzzy

import (
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/require"
)

var pool = []string{
	"aa",
	"zz",
	"story",
	"stry",
	"scenario",
	"scenario1",
	"scenario2",
	"feature",
	"actor",
	"rule",
	"cucumber",
	"gherkin",
	"add custom metadata",
	"sync custom metadata",
	"run custom metadata",
}

func Test_StringDistanceRanking(t *testing.T) {
	for _, test := range []string{
		"story",
		"scenario",
		"metadata",
		"s",
	} {
		t.Run(test, func(t *testing.T) {
			snaptest.Snapshot(t, Rank(test, pool))
		})
	}
}

func Test_Adjacency(t *testing.T) {
	require.True(t, Equivalent(Match{"", 0.82}, Match{"", 0.8}))
	require.False(t, Equivalent(Match{"", 0.52}, Match{"", 0.8}))
}
