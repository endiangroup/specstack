package specification

import (
	"fmt"
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/require"
)

func newMockScenario(t *testing.T, body string) *Scenario {
	fullBody := fmt.Sprintf("Feature:\n\n%s", body)
	spec := generateAndReadSpec(
		t, map[string]string{
			fmt.Sprintf("features/%s.feature", t.Name()): fullBody,
		},
	)
	return spec.Scenarios()[0]
}

func Test_AScenarioCanReturnItsBareString(t *testing.T) {
	raw := newMockScenario(
		t,
		`Scenario: Git not initialised for manual pull
    Given I have a project directory
    But I have not initialised git
    When I run "pull"
    Then I should see an error message informing me "initialise repository first"`,
	)

	snaptest.Snapshot(t, raw.bareString())
}

func Test_ScenarioRelated(t *testing.T) {
	startingScenario := newMockScenario(
		t,
		`Scenario: Git not initialised for manual pull
			Given I have a project directory
			But I have not initialised git
			When I run "pull"
			Then I should see an error message informing me "initialise repository first"`,
	)
	startingScenarioWithoutName := newMockScenario(
		t,
		`Scenario:
			Given I have a project directory
			But I have not initialised git
			When I run "pull"
			Then I should see an error message informing me "initialise repository first"`,
	)

	for _, test := range []struct {
		description string
		new         string
		old         *Scenario
		related     bool
	}{
		{
			description: "Identical scenarios",
			related:     true,
			new: `Scenario: Git not initialised for manual pull
					Given I have a project directory
					But I have not initialised git
					When I run "pull"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "Identical scenarios, small typos",
			related:     true,
			new: `Scenario: Git not initialised for manuall pull
					Given I have a project directory
					But I have not initialised git
					When I run "pull"
					Then I shoud see an error message informing me "initialise repository first"`,
		},
		{
			description: "Identical scenarios, different variables",
			related:     true,
			new: `Scenario: Git not initialised for manual pull
					Given I have a project directory
					But I have not initialised git
					When I run "command"
					Then I should see an error message informing me "some other value"`,
		},
		{
			description: "One line removed",
			related:     true,
			new: `Scenario: Git not initialised for manual pull
					Given I have a project directory
					But I have not initialised git
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "Two lines removed",
			related:     true,
			new: `Scenario: Git not initialised for manual pull
					Given I have a project directory
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "Three lines removed",
			related:     false,
			new: `Scenario: Git not initialised for manual pull
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "All lines removed, same name",
			related:     true,
			new:         `Scenario: Git not initialised for manual pull`,
		},
		{
			description: "Identical content, different name",
			related:     true,
			new: `Scenario: Some other name
					Given I have a project directory
					But I have not initialised git
					When I run "pull"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "Identical content, no name",
			related:     true,
			new: `Scenario: 
					Given I have a project directory
					But I have not initialised git
					When I run "pull"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "One line removed, different name",
			related:     true,
			new: `Scenario: Some other name
					Given I have a project directory
					But I have not initialised git
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "Two lines removed, different name",
			related:     false,
			new: `Scenario: Some other name
					Given I have a project directory
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "One line changed",
			related:     true,
			new: `Scenario: Git not initialised for manual pull
					Given Some other line
					But I have not initialised git
					When I run "pull"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "One line changed, another removed",
			related:     true,
			new: `Scenario: Git not initialised for manual pull
					Given Some other line
					But I have not initialised git
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "Two lines changed",
			related:     true,
			new: `Scenario: Git not initialised for manual pull
					Given Some other line
					But A totally different line
					When I run "pull"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "Three lines changed",
			related:     false,
			new: `Scenario: Git not initialised for manual pull
					Given Some other line
					But A totally different line
					When I run "something else"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "New scenario has a name, and the same content",
			related:     true,
			old:         startingScenarioWithoutName,
			new: `Scenario: Git not initialised for manual pull
					Given I have a project directory
					But I have not initialised git
					When I run "pull"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "New scenario has a name, one line removed",
			related:     true,
			old:         startingScenarioWithoutName,
			new: `Scenario: Git not initialised for manual pull
					Given I have a project directory
					When I run "pull"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "New scenario has a name, one line removed plus typos",
			related:     true,
			old:         startingScenarioWithoutName,
			new: `Scenario: Git not initialised for manuall pull
					Given I have a projec directory
					When I run "pull"
					Then I should see an error message informing me "initialise repository first"`,
		},
		{
			description: "New scenario has a name, two lines removed",
			related:     false,
			old:         startingScenarioWithoutName,
			new: `Scenario: Git not initialised for manual pull
					Given I have a project directory
					Then I should see an error message informing me "initialise repository first"`,
		},
	} {
		t.Run(fmt.Sprintf("input '%s'", test.description), func(t *testing.T) {
			old := startingScenario
			if test.old != nil {
				old = test.old
			}
			require.Equal(t, test.related, ScenarioRelated(old, newMockScenario(t, test.new)))
		})
	}
}
