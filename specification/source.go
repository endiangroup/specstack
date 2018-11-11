package specification

import (
	"sort"

	gherkin "github.com/DATA-DOG/godog/gherkin"
)

type Specification struct {
	FeatureFiles map[string]*gherkin.Feature
}

func NewSpecification() *Specification {
	return &Specification{
		FeatureFiles: make(map[string]*gherkin.Feature),
	}
}

// Stories fetches a list of features derived from loaded feature files.
// Features are returned in alphabetical order of the file name that contains
// them.
func (f *Specification) Stories() []*gherkin.Feature {
	stories := []*gherkin.Feature{}
	files := []string{}

	for file := range f.FeatureFiles {
		files = append(files, file)
	}

	sort.Strings(files)

	for _, index := range files {
		stories = append(stories, f.FeatureFiles[index])
	}

	return stories
}

// A Reader represents the input for a specification. The read method
// returns a Specification, zero or more warnings, and a fatal error.
type Reader interface {
	Read() (*Specification, []error, error)
}
