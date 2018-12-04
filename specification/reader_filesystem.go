package specification

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	// FIXME:	gherkin "github.com/cucumber/gherkin-go"
	// OR github.com/cucumber/cucumber/gherkin/go ?
	gherkin "github.com/DATA-DOG/godog/gherkin"
	"github.com/endiangroup/specstack/errors"
	"github.com/spf13/afero"
)

// File extensions for feature files
const (
	FileExtFeature = ".feature"
	FileExtStory   = ".story"
)

// A Filesystem represents a specification stored on a disk, memory, or other
// similar entity.
type Filesystem struct {
	Fs   afero.Fs
	Path string
}

// NewFilesystemReader creates a new Filesystem-based Source given an afero.Fs
// and a path. It scans the directory recursively, looking for .feature and .story
// files. It returns a Source, a list of warnings and an error.
func NewFilesystemReader(fs afero.Fs, path string) Reader {
	return &Filesystem{
		Fs:   fs,
		Path: path,
	}
}

// Read reads a specification from disk, returning the spec, any warnings, and
// possibly a fatal error.
func (f *Filesystem) Read() (*Specification, errors.Warnings, error) {
	spec := NewSpecification()
	spec.Source = f.Path
	warnings := errors.Warnings{}

	if err := afero.Walk(f.Fs, f.Path, f.featuresAndStoriesWalkFunc(spec, &warnings)); err != nil {
		return nil, warnings, fmt.Errorf("failed to read directory %s: %s", f.Path, err)
	}

	return spec, warnings, nil
}

func (f *Filesystem) ReadSource(s Sourcer) (io.Reader, error) {
	return f.Fs.Open(s.Source())
}

func (f *Filesystem) featuresAndStoriesWalkFunc(spec *Specification, warnings *errors.Warnings) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case FileExtFeature, FileExtStory:
			if err := f.addFeatureFile(spec, path); err != nil {
				*warnings = warnings.Append(err)
			}
		}

		return nil
	}
}

// addFeatureFile tries to parse a file in a given afero.Fs and adds it to the
// Filesystem state.
func (f *Filesystem) addFeatureFile(spec *Specification, path string) error {

	story, scenarios, err := f.parseFeatureFile(f.Fs, path)

	if err != nil {
		return err
	}

	spec.StorySources[path] = story
	spec.ScenarioSources[story] = scenarios

	return nil
}

func (f *Filesystem) parseFeatureFile(fs afero.Fs, path string) (*Story, []*Scenario, error) {
	content, err := afero.ReadFile(fs, path)

	if err != nil {
		return &Story{}, nil, fmt.Errorf("failed to read %s: %s", path, err)
	}

	buf := bytes.NewBuffer(content)
	feature, err := gherkin.ParseFeature(buf)

	if err != nil {
		return &Story{}, nil, fmt.Errorf("failed to parse %s: %s", path, err)
	}

	story := newStoryFromGherkinFeature(feature, path)
	scenarios := []*Scenario{}

	for _, s := range feature.ScenarioDefinitions {
		switch scenario := s.(type) {
		case *gherkin.Scenario:
			scenarios = append(scenarios, newScenarioFromGherkinScenario(scenario, story))
		default:
			return nil, nil, fmt.Errorf("Unhandled type %s", reflect.TypeOf(s))
		}
	}

	return story, scenarios, nil
}
