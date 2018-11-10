package specification

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
	FeatureFiles map[string]*gherkin.Feature
}

// NewFilesystem is a constructor for Filesystem. It allocates memory to the struct.
// Filesystems should generally be created with NewSourceFromFilesystem.
func NewFilesystem() *Filesystem {
	return &Filesystem{
		FeatureFiles: make(map[string]*gherkin.Feature),
	}
}

// NewSourceFromFilesystem creates a new Filesystem-based Source given an afero.Fs
// and a path. It scans the directory recursively, looking for .feature and .story
// files. It returns a Source, a list of warnings and an error.
func NewSourceFromFilesystem(fs afero.Fs, path string) (Source, []error, error) {
	reader := NewFilesystem()
	warnings := []error{}

	err := afero.Walk(fs, path, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		switch filepath.Ext(path) {
		case FileExtFeature, FileExtStory:

			if err := reader.AddFeatureFile(fs, path); err != nil {
				warnings = append(warnings, err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, warnings, fmt.Errorf("Failed to read directory %s: %s", path, err)
	}

	return reader, warnings, nil
}

// AddFeatureFile tries to parse a file in a given afero.Fs and adds it to the
// Filesystem state.
func (f *Filesystem) AddFeatureFile(fs afero.Fs, path string) error {

	feature, err := f.parseFeatureFile(fs, path)

	if err != nil {
		return err
	}

	f.FeatureFiles[path] = feature

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

// Stories fetches a list of features derived from loaded feature files.
// Features are returned in alphabetical order of the file name that contains
// them.
func (f *Filesystem) Stories() []*gherkin.Feature {
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
