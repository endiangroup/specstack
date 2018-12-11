package specification

import (
	"sort"
	"strings"

	"github.com/endiangroup/specstack/fuzzy"
)

type scenarioIndex struct {
	name     string
	scenario *Scenario
}

type scenarioIndexes []scenarioIndex

func (s scenarioIndexes) find(name string) *Scenario {
	for _, v := range s {
		if v.name == name {
			return v.scenario
		}
	}
	return nil
}

type FilterFunc func(string, []string) []string

type Filter struct {
	specification *Specification
	filterFuncs   []FilterFunc
	stories       []*Story
	scenarios     []*Scenario
}

func NewFilter(specification *Specification, filters ...FilterFunc) *Filter {
	filterFuncs := []FilterFunc{ClosestMatch}
	if len(filters) > 0 {
		filterFuncs = filters
	}

	return &Filter{
		specification: specification,
		filterFuncs:   filterFuncs,
	}
}

func (f *Filter) StoryQuery(term string) *Filter {
	f.stories = []*Story{}

	allStorySources := make(map[string]*Story)
	for k, v := range f.specification.StorySources {
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

	matches := finalSources
	for _, fn := range f.filterFuncs {
		matches = fn(term, matches)
	}

	for _, match := range matches {
		f.stories = append(f.stories, allStorySources[lookup[match]])
	}

	// Make sure no two matches are the same
	if len(f.stories) > 1 {
		for i := 1; i < len(f.stories); i++ {
			if f.stories[i-1] == f.stories[i] {
				f.stories = append(f.stories[:i-1], f.stories[i])
			}
		}
	}

	return f
}

func (f *Filter) ScenarioQuery(term string) *Filter {
	f.scenarios = []*Scenario{}
	allScenarios := scenarioIndexes{}
	uniqueNames := make(map[string]struct{})

	for _, s := range f.specification.Scenarios(f.stories...) {
		allScenarios = append(allScenarios, scenarioIndex{s.Name, s})
		uniqueNames[s.Name] = struct{}{}
	}

	pool := []string{}
	for k := range uniqueNames {
		pool = append(pool, k)
	}

	matches := pool
	for _, fn := range f.filterFuncs {
		matches = fn(term, matches)
	}

	for _, match := range matches {
		f.scenarios = append(f.scenarios, allScenarios.find(match))
	}

	return f
}

func (f *Filter) Stories() []*Story {
	return f.stories
}

func (f *Filter) Scenarios() []*Scenario {
	return f.scenarios
}

func (f *Filter) trimSource(input string) string {
	specSource := f.specification.Source + "/"
	trimmed := strings.TrimPrefix(input, specSource)
	trimmed = strings.TrimSuffix(trimmed, FileExtFeature)
	trimmed = strings.TrimSuffix(trimmed, FileExtStory)
	return trimmed
}

// ClosestMatch gets the best match from a pool of strings,
// omitting entries that fall below the fuzzy match threshold,
// and returning more than one result of they're roughly fuzzy-
// equivalent.
//
// This function is used as the default FilterFunc for Filters.
func ClosestMatch(term string, pool []string) []string {
	ranked := fuzzy.Rank(term, pool)

	if len(ranked) == 0 {
		return nil
	}

	if ranked[0].Negligible() {
		return nil
	}

	if len(ranked) > 1 && fuzzy.Equivalent(ranked[0], ranked[1]) {
		outputs := []string{ranked[0].Term, ranked[1].Term}
		sort.Strings(outputs)
		return outputs
	}

	return []string{ranked[0].Term}
}
