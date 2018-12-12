package specification

import (
	"fmt"
	"sort"
	"strconv"

	gherkin "github.com/DATA-DOG/godog/gherkin"
)

type Specification struct {
	Source          string
	StorySources    map[string]*Story
	ScenarioSources map[*Story][]*Scenario
}

func NewSpecification() *Specification {
	return &Specification{
		StorySources:    make(map[string]*Story),
		ScenarioSources: make(map[*Story][]*Scenario),
	}
}

type Story struct {
	*gherkin.Feature
	SourceIdentifier string
}

func newStoryFromGherkinFeature(feature *gherkin.Feature, source string) *Story {
	return &Story{
		Feature:          feature,
		SourceIdentifier: source,
	}
}

func (s *Story) Source() Source {
	return Source{SourceTypeFile, s.SourceIdentifier}
}

// Stories fetches a list of features derived from loaded feature files.
// Features are returned in alphabetical order of the file name that contains
// them.
func (f *Specification) Stories() []*Story {
	stories := []*Story{}
	sources := []string{}

	for file := range f.StorySources {
		sources = append(sources, file)
	}

	sort.Strings(sources)

	for _, index := range sources {
		stories = append(stories, f.StorySources[index])
	}

	return stories
}

// Scenarios fetches a complete list of scenarios from all loaded feature
// files. Scenarios are returned in the order they appear in their feature
// file, grouped by file name in alphabetical order.
func (s *Specification) Scenarios(stories ...*Story) []*Scenario {
	scenarios := []*Scenario{}

	if len(stories) == 0 {
		stories = s.Stories()
	}

	for _, story := range stories {
		scenarios = append(scenarios, s.ScenarioSources[story]...)
	}
	return scenarios
}

// FindStory performs a fuzzy match on the source (usually file name) and
// name of all known stories, then returns the closest match, if any. The base
// source (usually directory path) and any file extensions are omitted from the
// match. In the event of a tie (that is, two roughly equal matches) then an
// error is returned.
func (f *Specification) FindStory(input string) (*Story, error) {
	matches := NewFilter(f).Query(
		MapStories(ReduceClosestMatch(input)),
		DedupStories(),
	).Stories()

	if len(matches) == 0 {
		return nil, fmt.Errorf("no story matching %s", input)
	}

	if len(matches) > 1 {
		return nil, fmt.Errorf(
			"story name is ambiguous. The most similar story names are '%s' and '%s'",
			matches[0].Name,
			matches[1].Name,
		)
	}

	return matches[0], nil
}

// FindScenario performs a fuzzy match on the name of all scenarios
// in scope. The scope is either all scenarios, or only scenarios in
// the provided story name. In the event of a tie (that is, two roughly
// equal matches) an error is returned
func (s *Specification) FindScenario(query, storyName string) (*Scenario, error) {
	filter := NewFilter(s)
	if storyName != "" {
		filter.Query(
			MapStories(ReduceClosestMatch(storyName)),
			DedupStories(),
		)
	}

	if val, err := strconv.Atoi(query); err == nil {
		filter.Query(MapScenarioIndex(val))
	} else {
		filter.Query(MapScenarios(ReduceClosestMatch(query)))
	}

	matches := filter.Scenarios()

	if len(matches) == 0 {
		return nil, fmt.Errorf("no scenario matching %s", query)
	}

	if len(matches) > 1 {
		m0, m1 := matches[0], matches[1]
		name0, name1 := m0.Name, m1.Name

		if m0.Story != m1.Story {
			name0 = fmt.Sprintf("%s/%s", m0.Story.Name, m0.Name)
			name1 = fmt.Sprintf("%s/%s", m1.Story.Name, m1.Name)
		}

		return nil, fmt.Errorf(
			"scenario query is ambiguous. The most similar scenario names are '%s' and '%s'",
			name0,
			name1,
		)
	}

	return matches[0], nil
}
