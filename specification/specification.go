package specification

import (
	"fmt"
	"sort"
	"strings"

	gherkin "github.com/DATA-DOG/godog/gherkin"
	"github.com/endiangroup/specstack/fuzzy"
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

func (s *Story) Source() string {
	return s.SourceIdentifier
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
func (s *Specification) Scenarios() []*Scenario {
	scenarios := []*Scenario{}
	for _, story := range s.Stories() {
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
	allStorySources := make(map[string]*Story)
	for k, v := range f.StorySources {
		allStorySources[k] = v
		allStorySources[v.Name] = v
	}

	sources := []string{}
	for file := range allStorySources {
		sources = append(sources, file)
	}

	lookup := make(map[string]string)
	for _, source := range sources {
		lookup[f.trimSource(source)] = source
	}

	finalSources := []string{}
	for fs := range lookup {
		finalSources = append(finalSources, fs)
	}

	matches := f.closetMatch(input, finalSources)

	if matches == nil {
		return nil, fmt.Errorf("no story matching %s", input)
	}

	if len(matches) > 1 {
		// The matches may point to the same underlying object
		if a, b := allStorySources[lookup[matches[0]]], allStorySources[lookup[matches[1]]]; a == b {
			return a, nil
		}

		return nil, fmt.Errorf(
			"story name is ambiguous. The most similar story names are '%s' and '%s'",
			matches[0],
			matches[1],
		)
	}

	return allStorySources[lookup[matches[0]]], nil
}

func (s *Specification) FindScenario(query string) (*Scenario, error) {
	return nil, fmt.Errorf("TODO: implement and test FindScenario")
}

func (f *Specification) closetMatch(term string, pool []string) []string {
	ranked := fuzzy.Rank(term, pool)

	if len(ranked) == 0 {
		return nil
	}

	if ranked[0].Negligible() {
		return nil
	}

	if len(ranked) > 1 && fuzzy.Adjacent(ranked[0], ranked[1]) {
		outputs := []string{ranked[0].Term, ranked[1].Term}
		sort.Strings(outputs)
		return outputs
	}

	return []string{ranked[0].Term}
}

func (f *Specification) trimSource(input string) string {
	specSource := f.Source + "/"
	trimmed := strings.TrimPrefix(input, specSource)
	trimmed = strings.TrimSuffix(trimmed, FileExtFeature)
	trimmed = strings.TrimSuffix(trimmed, FileExtStory)
	return trimmed
}
