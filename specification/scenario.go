package specification

import (
	"strings"

	gherkin "github.com/DATA-DOG/godog/gherkin"
	"github.com/endiangroup/specstack/fuzzy"
)

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

func (s *Scenario) Source() string {
	return "FIXME: Do we need this?"
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

func ScenarioDistance(a, b *Scenario) float64 {
	if a.Name == "" || b.Name == "" {
		return fuzzy.Strcmp(a.bareStringStepsOnly(), b.bareStringStepsOnly())
	} else if len(a.Steps) == 0 || len(b.Steps) == 0 {
		return fuzzy.Strcmp(a.Name, b.Name)
	}
	return fuzzy.Strcmp(a.bareString(), b.bareString())
}

func ScenarioRelated(a, b *Scenario) bool {
	return ScenarioDistance(a, b) >= fuzzy.DistanceThreshold
}
