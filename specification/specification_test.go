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

  Scenario: First
	Then this will work
`
	mockFeatureC = `Feature: completely different
  In order to test application behavior
  As a test suite
  I need to be able to run features

  Scenario: Second
	Then this will work
`
	mockFeatureD = `Feature: manage metadata
  In order to test application behavior
  As a test suite
  I need to be able to run features

  Scenario: Third
	Then this will work
`
	mockFeatureE = `Feature: ABC
  In order to test application behavior
  As a test suite
  I need to be able to run features

  Scenario: Fourth
	Then this will work
`
	mockFeatureF = `Feature: Very similar1

  Scenario: Very similar1
	Then this will work
`
	mockFeatureG = `Feature: Very similar2

  Scenario: Very similar2
	Then this will work
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
			"features/similar1.feature":      mockFeatureF,
			"features/similar2.feature":      mockFeatureG,
			"features/BBC.feature":           mockFeatureE,
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
		{term: "similar", err: fmt.Errorf("story name is ambiguous. The most similar story names are 'Very similar1' and 'Very similar2'")}, //nolint:lll
		// In this case, there are two equal matches (the file name
		// and the story title) which point to the same thing, so there
		// should be no error
		{term: "C"},
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

func Test_ASpecificationCanAddressScenarios(t *testing.T) {

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
		term  string
		story string
		err   error
	}{
		{term: "first"},
		{term: "zzz", err: fmt.Errorf("no scenario matching zzz")},
		{term: "similar", story: "similar1"},
		{term: "similar", err: fmt.Errorf("scenario query is ambiguous. The most similar scenario names are 'Very similar1/Very similar1' and 'Very similar2/Very similar2'")}, //nolint:lll
	} {
		t.Run(test.term, func(t *testing.T) {
			story, err := spec.FindScenario(test.term, test.story)

			if test.err == nil {
				require.Nil(t, err)
				//snaptest.Snapshot(t, story)
			} else {
				require.Equal(t, test.err, err)
				require.Nil(t, story)
			}
		})
	}
}
