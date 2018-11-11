package specification

import (
	"fmt"
	"sort"
	"strings"

	gherkin "github.com/DATA-DOG/godog/gherkin"
	"github.com/schollz/closestmatch"
)

type Specification struct {
	Source       string
	StorySources map[string]*Story
}

func NewSpecification() *Specification {
	return &Specification{
		StorySources: make(map[string]*Story),
	}
}

type Story struct {
	*gherkin.Feature
	Source string
}

func newStoryFromGherkinFeature(feature *gherkin.Feature, source string) *Story {
	return &Story{
		Feature: feature,
		Source:  source,
	}
}

// Stories fetches a list of features derived from loaded feature files.
// Features are returned in alphabetical order of the file name that contains
// them.
func (f *Specification) Stories() []*Story {
	stories := []*Story{}
	sources := []string{}

	for file := range f.StorySources {
		sources = append(sources, file)
	}

	sort.Strings(sources)

	for _, index := range sources {
		stories = append(stories, f.StorySources[index])
	}

	return stories
}

// FindStory performs a fuzzy match on the source (usually file name) and
// name of all known stories, then returns the closest match. The base source
// (usually directory path) and any file extensions are ommitted from the
// match. In the event of a tie (that is, two equal matches) the story is chosen
// on its alphabetical primacy.
func (f *Specification) FindStory(input string) (*Story, error) {
	allStorySources := make(map[string]*Story)

	for k, v := range f.StorySources {
		allStorySources[k] = v
		allStorySources[v.Name] = v
	}

	sources := []string{}
	for file := range allStorySources {
		sources = append(sources, file)
	}

	lookup := make(map[string]string)
	for _, source := range sources {
		lookup[f.trimSource(source)] = source
	}

	finalSources := []string{}
	for fs := range lookup {
		finalSources = append(finalSources, fs)
	}

	sort.Strings(finalSources)
	cm := closestmatch.New(finalSources, []int{2, 3})
	match := cm.Closest(input)

	if match == "" {
		return nil, fmt.Errorf("no story matching %s", input)
	}

	return allStorySources[lookup[match]], nil
}

func (f *Specification) trimSource(input string) string {
	specSource := f.Source + "/"
	trimmed := strings.TrimPrefix(input, specSource)
	trimmed = strings.TrimSuffix(trimmed, FileExtFeature)
	trimmed = strings.TrimSuffix(trimmed, FileExtStory)
	return trimmed
}
