package specification

import gherkin "github.com/DATA-DOG/godog/gherkin"

// FIXME: import gherkin "github.com/cucumber/gherkin-go"

type Reader interface {
	Stories() []gherkin.Feature
}
