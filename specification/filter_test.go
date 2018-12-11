package specification

import (
	"testing"

	"github.com/endiangroup/snaptest"
)

func Test_AFilterCanAddressStories(t *testing.T) {

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
			filter := NewFilter(spec)
			filter.StoryQuery(test.term)
			snaptest.Snapshot(t, filter.Stories())
		})
	}
}

func Test_AFilterCanAddressScenarios(t *testing.T) {

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
			filter := NewFilter(spec)
			filter.ScenarioQuery(test.term)
			snaptest.Snapshot(t, filter.Scenarios())
		})
	}
}

func Test_AFilterCanAddressAScenarioInAStory(t *testing.T) {
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

	filter := NewFilter(spec)
	filter.StoryQuery("similar1")
	filter.ScenarioQuery("similar")
	snaptest.Snapshot(t, filter.Scenarios())
}
