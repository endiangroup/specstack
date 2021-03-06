package specification

import (
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/require"
)

func Test_AQueryCanAddressStories(t *testing.T) {

	spec := generateAndReadSpec(t,
		map[string]string{
			"features/set_up_repo.feature":   mockFeatureA,
			"features/update_config.feature": mockFeatureB,
			"features/create_config.feature": mockFeatureC,
			"features/add_metadata.feature":  mockFeatureD,
			"features/similar1.feature":      mockFeatureF,
			"features/similar2.feature":      mockFeatureG,
			"features/BBC.feature":           mockFeatureE,
		},
	)

	for _, test := range []struct {
		term string
	}{
		{term: "set_up_repo"},
		{term: "different"},
		{term: "similar"},
		{term: "zzz"},
	} {
		t.Run(test.term, func(t *testing.T) {
			filter := NewQuery(spec)
			filter.MapReduce(MapStories(ReduceClosestMatch(test.term)), MapUniqueStories())
			snaptest.Snapshot(t, filter.Stories())
		})
	}
}

func Test_AQueryCanAddressScenarios(t *testing.T) {

	spec := generateAndReadSpec(t,
		map[string]string{
			"features/b.feature": mockFeatureB,
			"features/c.feature": mockFeatureC,
			"features/d.feature": mockFeatureD,
			"features/e.feature": mockFeatureE,
			"features/f.feature": mockFeatureF,
			"features/g.feature": mockFeatureG,
		},
	)

	for _, test := range []struct {
		term string
	}{
		{term: "first"},
		{term: "similar"},
		{term: "zzz"},
	} {
		t.Run(test.term, func(t *testing.T) {
			filter := NewQuery(spec)
			filter.MapReduce(MapScenarios(ReduceClosestMatch(test.term)))
			snaptest.Snapshot(t, filter.Scenarios())
		})
	}
}

func Test_AQueryCanAddressAScenarioInAStory(t *testing.T) {
	spec := generateAndReadSpec(t,
		map[string]string{
			"features/set_up_repo.feature":   mockFeatureA,
			"features/update_config.feature": mockFeatureB,
			"features/create_config.feature": mockFeatureC,
			"features/add_metadata.feature":  mockFeatureD,
			"features/similar1.feature":      mockFeatureF,
			"features/similar2.feature":      mockFeatureG,
			"features/BBC.feature":           mockFeatureE,
		},
	)

	filter := NewQuery(spec).MapReduce(
		MapStories(ReduceClosestMatch("similiar1")),
		MapUniqueStories(),
		MapScenarios(ReduceClosestMatch("similiar")),
	)
	snaptest.Snapshot(t, filter.Scenarios())
}

func Test_AQueryCanIndexAScenarioInAStory(t *testing.T) {
	spec := generateAndReadSpec(t,
		map[string]string{
			"features/two_scenarios.feature": mockFeatureH,
		},
	)

	require.Equal(t, "Scenario A",
		NewQuery(spec).MapReduce(
			MapStories(ReduceClosestMatch("two_scenarios")),
			MapScenarioIndex(1),
		).Scenarios()[0].Name,
	)

	require.Equal(t, "Scenario B",
		NewQuery(spec).MapReduce(
			MapStories(ReduceClosestMatch("two_scenarios")),
			MapScenarioIndex(2),
		).Scenarios()[0].Name,
	)

	require.Empty(t,
		NewQuery(spec).MapReduce(
			MapStories(ReduceClosestMatch("two_scenarios")),
			MapScenarioIndex(3),
		).Scenarios(),
	)
}

func Test_AQueryCanAddressAScenarioByFuncFilter(t *testing.T) {
	spec := generateAndReadSpec(t,
		map[string]string{
			"features/set_up_repo.feature":   mockFeatureA,
			"features/update_config.feature": mockFeatureB,
			"features/create_config.feature": mockFeatureC,
			"features/add_metadata.feature":  mockFeatureD,
			"features/similar1.feature":      mockFeatureF,
			"features/similar2.feature":      mockFeatureG,
			"features/BBC.feature":           mockFeatureE,
		},
	)

	filter := NewQuery(spec).MapReduce(
		MapScenarios(),
		MapScenarioMatchFunc(func(s *Scenario) bool {
			return s.Name == "Second"
		}),
	)
	snaptest.Snapshot(t, filter.Scenarios())
}
