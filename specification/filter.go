package specification

import (
	"sort"
	"strings"

	"github.com/endiangroup/specstack/fuzzy"
)

type FilterMapFunc func(*Filter)
type FilterReduceFunc func([]string) []string

type Filter struct {
	specification *Specification
	stories       []*Story
	scenarios     []*Scenario
}

func NewFilter(specification *Specification) *Filter {
	return &Filter{
		specification: specification,
	}
}

func (f *Filter) MapReduce(fns ...FilterMapFunc) *Filter {
	for _, fn := range fns {
		fn(f)
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

func (f *Filter) applyReduceFns(input []string, fns []FilterReduceFunc) []string {
	matches := input
	for _, fn := range fns {
		matches = fn(matches)
	}
	return matches
}

func (f *Filter) storySources() (map[string]*Story, []string) {
	allStorySources := make(map[string]*Story)
	for k, v := range f.specification.StorySources {
		allStorySources[k] = v
		allStorySources[v.Name] = v
	}

	sources := []string{}
	for file := range allStorySources {
		sources = append(sources, file)
	}

	return allStorySources, sources
}

func MapStories(filters ...FilterReduceFunc) FilterMapFunc {
	return func(f *Filter) {
		f.stories = []*Story{}
		allStorySources, sources := f.storySources()

		lookup := make(map[string]string)
		for _, source := range sources {
			lookup[f.trimSource(source)] = source
		}

		finalSources := []string{}
		for fs := range lookup {
			finalSources = append(finalSources, fs)
		}

		matches := f.applyReduceFns(finalSources, filters)
		for _, match := range matches {
			f.stories = append(f.stories, allStorySources[lookup[match]])
		}
	}
}

func MapUniqueStories() FilterMapFunc {
	return func(f *Filter) {
		if len(f.stories) <= 1 {
			return
		}
		for i := 1; i < len(f.stories); i++ {
			if f.stories[i-1] == f.stories[i] {
				f.stories = append(f.stories[:i-1], f.stories[i])
			}
		}
	}
}

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

func MapScenarios(filters ...FilterReduceFunc) FilterMapFunc {
	return func(f *Filter) {
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

		matches := f.applyReduceFns(pool, filters)
		for _, match := range matches {
			f.scenarios = append(f.scenarios, allScenarios.find(match))
		}
	}
}

func MapScenarioIndex(index int) FilterMapFunc {
	return func(f *Filter) {
		f.scenarios = []*Scenario{}
		if len(f.stories) != 1 {
			return
		}

		scenarios := f.specification.Scenarios(f.stories[0])
		if index > len(scenarios) {
			return
		}

		f.scenarios = []*Scenario{scenarios[index-1]}
	}
}

func ReduceClosestMatch(term string) FilterReduceFunc {
	return func(pool []string) []string {
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
}
