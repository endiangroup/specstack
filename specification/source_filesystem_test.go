package specification

import (
	"fmt"
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

const (
	mockFeatureA = `Feature: run features
  In order to test application behavior
  As a test suite
  I need to be able to run features

  Scenario: should run a normal feature
    Given a feature "normal.feature" file:
      """
      Feature: normal feature

        Scenario: parse a scenario
          Given a feature path "features/load.feature:6"
          When I parse features
          Then I should have 1 scenario registered
      """
    When I run feature suite
    Then the suite should have passed
    And the following steps should be passed:
`
)

func newSpecificationFs(t *testing.T, files map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()

	for path, content := range files {
		file, err := fs.Create(path)
		require.Nil(t, err)
		_, err = file.WriteString(content)
		require.Nil(t, err)
	}

	return fs
}

func Test_AFilesystemReaderCanReadAFeatureFileFromDisk(t *testing.T) {

	for _, test := range []struct {
		description string
		fileContent map[string]string
		inputPath   string
		err         error
	}{
		{
			description: "Happy path",
			fileContent: map[string]string{"features/a.feature": mockFeatureA},
			inputPath:   "features/a.feature",
		},
		{
			description: "Sad path: file doesn't exist",
			fileContent: map[string]string{},
			inputPath:   "features/a.feature",
			err:         fmt.Errorf("Failed to read features/a.feature: open features/a.feature: file does not exist"),
		},
		{
			description: "Sad path: file content invalid",
			fileContent: map[string]string{"features/a.feature": "--invalid--"},
			inputPath:   "features/a.feature",
			err:         fmt.Errorf("Failed to parse features/a.feature: Parser errors:\n(1:1): expected: #Language, #TagLine, #FeatureLine, #Comment, #Empty, got '--invalid--'\n(2:0): unexpected end of file, expected: #Language, #TagLine, #FeatureLine, #Comment, #Empty"),
		},
	} {
		t.Run(fmt.Sprintf("input '%s'", test.description), func(t *testing.T) {
			spec := NewSpecification()
			reader := Filesystem{
				Fs: newSpecificationFs(t, test.fileContent),
			}
			err := reader.addFeatureFile(spec, test.inputPath)

			if test.err == nil {
				require.Nil(t, err)
				snaptest.Snapshot(t, spec)
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}

func Test_AFilesystemReaderCanReadASpecificationFromDisk(t *testing.T) {

	for _, test := range []struct {
		description string
		fileContent map[string]string
		inputDir    string
		warnings    []error
		err         error
	}{
		{
			description: "Happy path: easy",
			fileContent: map[string]string{
				"features/a.feature": mockFeatureA,
				"features/b.feature": mockFeatureA,
			},
			inputDir: "features",
			warnings: []error{},
		},
		{
			description: "Happy path: non-feature files",
			fileContent: map[string]string{
				"features/a.feature":    mockFeatureA,
				"features/b.notfeature": "Not a feature file",
			},
			inputDir: "features",
			warnings: []error{},
		},
		{
			description: "Happy path: warnings",
			fileContent: map[string]string{
				"features/a.feature": mockFeatureA,
				"features/b.feature": "--invalid--",
			},
			inputDir: "features",
			warnings: []error{
				fmt.Errorf("Failed to parse features/b.feature: Parser errors:\n(1:1): expected: #Language, #TagLine, #FeatureLine, #Comment, #Empty, got '--invalid--'\n(2:0): unexpected end of file, expected: #Language, #TagLine, #FeatureLine, #Comment, #Empty"),
			},
		},
		{
			description: "Sad path: no dir",
			fileContent: map[string]string{
				"features/a.feature": mockFeatureA,
			},
			inputDir: "notfeatures",
			warnings: []error{},
			err:      fmt.Errorf("Failed to read directory notfeatures: open notfeatures: file does not exist"),
		},
	} {
		t.Run(fmt.Sprintf("input '%s'", test.description), func(t *testing.T) {
			fs := newSpecificationFs(t, test.fileContent)
			reader := NewFilesystemReader(fs, test.inputDir)
			spec, warnings, err := reader.Read()

			if test.err == nil {
				require.Nil(t, err)
				require.Equal(t, test.warnings, warnings)
				snaptest.Snapshot(t, spec)
			} else {
				require.Equal(t, test.err, err)
			}
		})
	}
}

func Test_ASpecificationCanGetAListOfStories(t *testing.T) {
	fs := newSpecificationFs(t,
		map[string]string{
			"features/a.feature": mockFeatureA,
			"features/b.feature": mockFeatureA,
		},
	)
	reader := NewFilesystemReader(fs, "features")
	require.NotNil(t, reader)

	spec, warnings, err := reader.Read()
	require.Nil(t, err)
	require.Len(t, warnings, 0)
	snaptest.Snapshot(t, spec.Stories())
}
