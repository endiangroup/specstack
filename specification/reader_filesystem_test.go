package specification

import (
	"fmt"
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
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
			err:         fmt.Errorf("Failed to parse features/a.feature: Parser errors:\n(1:1): expected: #Language, #TagLine, #FeatureLine, #Comment, #Empty, got '--invalid--'\n(2:0): unexpected end of file, expected: #Language, #TagLine, #FeatureLine, #Comment, #Empty"), //nolint:lll
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
				fmt.Errorf("Failed to parse features/b.feature: Parser errors:\n(1:1): expected: #Language, #TagLine, #FeatureLine, #Comment, #Empty, got '--invalid--'\n(2:0): unexpected end of file, expected: #Language, #TagLine, #FeatureLine, #Comment, #Empty"), //nolint:lll
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
