package specification

import (
	"math"
	"strings"

	gherkin "github.com/DATA-DOG/godog/gherkin"
	"github.com/antzucaro/matchr"
)

const distanceThreshold = 0.75

type Scenario struct {
	*gherkin.Scenario
	Story *Story
}

func newScenarioFromGherkinScenario(scenario *gherkin.Scenario, story *Story) *Scenario {
	return &Scenario{
		Scenario: scenario,
		Story:    story,
	}
}

func (s *Scenario) bareLines() []string {
	output := []string{s.Name}
	for _, step := range s.Steps {
		output = append(output, step.Text)
	}
	return output
}

func (s *Scenario) bareString() string {
	return strings.Join(s.bareLines(), "\n")
}

func (s *Scenario) bareStringStepsOnly() string {
	return strings.Join(s.bareLines()[1:], "\n")
}

func levstrcmp(a, b string) float64 {
	maxLen := math.Max(float64(len(a)), float64(len(b)))
	levDist := matchr.DamerauLevenshtein(a, b)
	return (1 - (float64(levDist) / maxLen))
}

func ScenarioDistance(a, b *Scenario) float64 {
	if a.Name == "" || b.Name == "" {
		return levstrcmp(a.bareStringStepsOnly(), b.bareStringStepsOnly())
	} else if len(a.Steps) == 0 || len(b.Steps) == 0 {
		return levstrcmp(a.Name, b.Name)
	}
	return levstrcmp(a.bareString(), b.bareString())
}

func ScenarioRelated(a, b *Scenario) bool {
	return ScenarioDistance(a, b) >= distanceThreshold
}
