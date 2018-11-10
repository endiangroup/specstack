package specification

import gherkin "github.com/DATA-DOG/godog/gherkin"

// FIXME: import gherkin "github.com/cucumber/gherkin-go"

// A Source represents the input for a specification.
type Source interface {
	Stories() []*gherkin.Feature
}
