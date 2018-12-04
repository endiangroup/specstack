package specification

import (
	"fmt"
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/require"
)

const (
	mockFeatureA = `Feature: run features
  In order to test application behavior
  As a test suite
  I need to be able to run features

  Scenario: should run a normal feature
    Given a feature "normal.feature" file:
      """
      Feature: normal feature

        Scenario: parse a scenario
          Given a feature path "features/load.feature:6"
          When I parse features
          Then I should have 1 scenario registered
      """
    When I run feature suite
    Then the suite should have passed
    And the following steps should be passed:
`
	mockFeatureB = `Feature: Search for me
  In order to test application behavior
  As a test suite
  I need to be able to run features
`
	mockFeatureC = `Feature: completely different
  In order to test application behavior
  As a test suite
  I need to be able to run features
`
	mockFeatureD = `Feature: manage metadata
  In order to test application behavior
  As a test suite
  I need to be able to run features
`
)

func generateAndReadSpec(t *testing.T, files map[string]string) *Specification {
	fs := newSpecificationFs(t, files)
	reader := NewFilesystemReader(fs, "features")
	require.NotNil(t, reader)

	spec, warnings, err := reader.Read()
	require.Nil(t, err)
	require.Len(t, warnings, 0)

	return spec
}

func Test_ASpecificationCanGetAListOfStories(t *testing.T) {
	spec := generateAndReadSpec(t,
		map[string]string{
			"features/a.feature": mockFeatureA,
			"features/b.feature": mockFeatureA,
		},
	)
	snaptest.Snapshot(t, spec.Stories())
}

func Test_ASpecificationCanAddressStories(t *testing.T) {

	spec := generateAndReadSpec(t,
		map[string]string{
			"features/set_up_repo.feature":   mockFeatureA,
			"features/update_config.feature": mockFeatureB,
			"features/create_config.feature": mockFeatureC,
			"features/add_metadata.feature":  mockFeatureD,
		},
	)

	for _, test := range []struct {
		term string
		err  error
	}{
		{term: "set_up_repo"},
		{term: "upconfig"},
		{term: "crcon"},
		{term: "Search for me"},
		{term: "different"},
		{term: "managemeta"},
		{term: "zzz", err: fmt.Errorf("no story matching zzz")},
	} {
		t.Run(test.term, func(t *testing.T) {
			story, err := spec.FindStory(test.term)

			if test.err == nil {
				require.Nil(t, err)
				snaptest.Snapshot(t, story)
			} else {
				require.Equal(t, test.err, err)
				require.Nil(t, story)
			}
		})
	}
}

func Test_ASpecificationCanGetAListOfScenarios(t *testing.T) {
	spec := generateAndReadSpec(t,
		map[string]string{
			"features/a.feature": mockFeatureA,
			"features/b.feature": mockFeatureA,
		},
	)
	snaptest.Snapshot(t, spec.Scenarios())
}
