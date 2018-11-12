package specification

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

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

	err := afero.Walk(f.Fs, f.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case FileExtFeature, FileExtStory:
			if err := f.addFeatureFile(spec, path); err != nil {
				warnings = append(warnings, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, warnings, fmt.Errorf("Failed to read directory %s: %s", f.Path, err)
	}

	return spec, warnings, nil
}

// addFeatureFile tries to parse a file in a given afero.Fs and adds it to the
// Filesystem state.
func (f *Filesystem) addFeatureFile(spec *Specification, path string) error {

	story, err := f.parseFeatureFile(f.Fs, path)

	if err != nil {
		return err
	}

	spec.StorySources[path] = story

	return nil
}

func (f *Filesystem) parseFeatureFile(fs afero.Fs, path string) (*Story, error) {
	content, err := afero.ReadFile(fs, path)

	if err != nil {
		// FIXME! Custom error for warnings?
		return &Story{}, fmt.Errorf("Failed to read %s: %s", path, err)
	}

	buf := bytes.NewBuffer(content)
	feature, err := gherkin.ParseFeature(buf)

	if err != nil {
		// FIXME! Custom error for warnings?
		return &Story{}, fmt.Errorf("Failed to parse %s: %s", path, err)
	}

	return newStoryFromGherkinFeature(feature, path), nil
}
