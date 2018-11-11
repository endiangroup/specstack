package specification

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	// FIXME:	gherkin "github.com/cucumber/gherkin-go"
	// OR github.com/cucumber/cucumber/gherkin/go ?
	gherkin "github.com/DATA-DOG/godog/gherkin"
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

func (f *Filesystem) Read() (*Specification, []error, error) {
	reader := NewSpecification()
	warnings := []error{}

	err := afero.Walk(f.Fs, f.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case FileExtFeature, FileExtStory:

			if err := f.AddFeatureFile(reader, path); err != nil {
				warnings = append(warnings, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, warnings, fmt.Errorf("Failed to read directory %s: %s", f.Path, err)
	}

	return reader, warnings, nil
}

// AddFeatureFile tries to parse a file in a given afero.Fs and adds it to the
// Filesystem state.
func (f *Filesystem) AddFeatureFile(spec *Specification, path string) error {

	feature, err := f.parseFeatureFile(f.Fs, path)

	if err != nil {
		return err
	}

	spec.FeatureFiles[path] = feature

	return nil
}

func (f *Filesystem) parseFeatureFile(fs afero.Fs, path string) (*gherkin.Feature, error) {
	content, err := afero.ReadFile(fs, path)

	if err != nil {
		// FIXME! Custom error for warnings?
		return &gherkin.Feature{}, fmt.Errorf("Failed to read %s: %s", path, err)
	}

	buf := bytes.NewBuffer(content)
	feature, err := gherkin.ParseFeature(buf)

	if err != nil {
		// FIXME! Custom error for warnings?
		return &gherkin.Feature{}, fmt.Errorf("Failed to parse %s: %s", path, err)
	}

	return feature, nil
}
