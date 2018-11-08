package specification

import (
	"bytes"
	"fmt"

	// FIXME:	gherkin "github.com/cucumber/gherkin-go"
	// OR github.com/cucumber/cucumber/gherkin/go ?
	gherkin "github.com/DATA-DOG/godog/gherkin"
	"github.com/spf13/afero"
)

type filesystemReader struct {
	featureFiles map[string]*gherkin.Feature
}

func (f *filesystemReader) Stories() []gherkin.Feature {
	// TODO
	return []gherkin.Feature{}
}

func NewFilesystemReader(fs afero.Fs) (Reader, error) {
	reader := &filesystemReader{
		featureFiles: make(map[string]*gherkin.Feature),
	}

	// FIXME! Paths. Maybe use readall? .story files?
	matches, err := afero.Glob(fs, "*.feature")

	if err != nil {
		return nil, err
	}

	for _, match := range matches {

		content, err := afero.ReadFile(fs, match)

		if err != nil {
			// FIXME! Custom error
			return nil, fmt.Errorf("Failed to read %s: %s", match, err)
		}

		buf := bytes.NewBuffer(content)
		feature, err := gherkin.ParseFeature(buf)

		if err != nil {
			// FIXME! Custom error
			return nil, fmt.Errorf("Failed to parse %s: %s", match, err)
		}

		reader.featureFiles[match] = feature
	}

	return reader, nil
}
