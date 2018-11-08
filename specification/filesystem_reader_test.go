package specification

import (
	"fmt"
	"testing"

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

func Test_AFilesystemReaderCanReadSpecificationFromDisk(t *testing.T) {
	fs := newSpecificationFs(t, map[string]string{
		"a.feature": `Feature: run features
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
`})

	reader, err := NewFilesystemReader(fs)
	require.Nil(t, err)
	fmt.Println(reader)
}
