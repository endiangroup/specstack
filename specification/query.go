package specification

import (
	"sort"
	"strings"

	"github.com/endiangroup/specstack/fuzzy"
)

type QueryMapFunc func(*Query)
type QueryReduceFunc func([]string) []string

type Query struct {
	specification *Specification
	stories       []*Story
	scenarios     []*Scenario
}

func NewQuery(specification *Specification) *Query {
	return &Query{
		specification: specification,
	}
}

func (q *Query) MapReduce(fns ...QueryMapFunc) *Query {
	for _, fn := range fns {
		fn(q)
	}
	return q
}

func (q *Query) Stories() []*Story {
	return q.stories
}

func (q *Query) Scenarios() []*Scenario {
	return q.scenarios
}

func (q *Query) trimSource(input string) string {
	specSource := q.specification.Source + "/"
	trimmed := strings.TrimPrefix(input, specSource)
	trimmed = strings.TrimSuffix(trimmed, FileExtFeature)
	trimmed = strings.TrimSuffix(trimmed, FileExtStory)
	return trimmed
}

func (q *Query) applyReduceFns(input []string, fns []QueryReduceFunc) []string {
	matches := input
	for _, fn := range fns {
		matches = fn(matches)
	}
	return matches
}

func (q *Query) storySources() (map[string]*Story, []string) {
	allStorySources := make(map[string]*Story)
	for k, v := range q.specification.StorySources {
		allStorySources[k] = v
		allStorySources[v.Name] = v
	}

	sources := []string{}
	for file := range allStorySources {
		sources = append(sources, file)
	}

	return allStorySources, sources
}

func MapStories(filters ...QueryReduceFunc) QueryMapFunc {
	return func(q *Query) {
		q.stories = []*Story{}
		allStorySources, sources := q.storySources()

		lookup := make(map[string]string)
		for _, source := range sources {
			lookup[q.trimSource(source)] = source
		}

		finalSources := []string{}
		for fs := range lookup {
			finalSources = append(finalSources, fs)
		}

		matches := q.applyReduceFns(finalSources, filters)
		for _, match := range matches {
			q.stories = append(q.stories, allStorySources[lookup[match]])
		}
	}
}

func MapUniqueStories() QueryMapFunc {
	return func(q *Query) {
		if len(q.stories) <= 1 {
			return
		}
		for i := 1; i < len(q.stories); i++ {
			if q.stories[i-1] == q.stories[i] {
				q.stories = append(q.stories[:i-1], q.stories[i])
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

func MapScenarios(filters ...QueryReduceFunc) QueryMapFunc {
	return func(q *Query) {
		q.scenarios = []*Scenario{}
		allScenarios := scenarioIndexes{}
		uniqueNames := make(map[string]struct{})

		for _, s := range q.specification.Scenarios(q.stories...) {
			allScenarios = append(allScenarios, scenarioIndex{s.Name, s})
			uniqueNames[s.Name] = struct{}{}
		}

		pool := []string{}
		for k := range uniqueNames {
			pool = append(pool, k)
		}

		matches := q.applyReduceFns(pool, filters)
		for _, match := range matches {
			q.scenarios = append(q.scenarios, allScenarios.find(match))
		}
	}
}

func MapScenarioIndex(index int) QueryMapFunc {
	return func(q *Query) {
		q.scenarios = []*Scenario{}
		if len(q.stories) != 1 {
			return
		}

		scenarios := q.specification.Scenarios(q.stories[0])
		if index > len(scenarios) {
			return
		}

		q.scenarios = []*Scenario{scenarios[index-1]}
	}
}

func ReduceClosestMatch(term string) QueryReduceFunc {
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

func ReduceMax(num int) QueryReduceFunc {
	return func(pool []string) []string {
		if len(pool) > num {
			pool = pool[:num]
		}
		return pool
	}
}
