package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func tempDirectory(t *testing.T) (path string, shutdown func()) {

	dir, err := ioutil.TempDir("", "specstack-test")
	require.Nil(t, err)
	require.NotEmpty(t, dir)

	return dir, func() {
		require.Nil(t, os.RemoveAll(dir))
	}
}

func initialisedGitRepoDir(t *testing.T) (path string, r *repositoryGit, shutdown func()) {

	dir, shutdown := tempDirectory(t)
	repo := NewGit(dir, "unittest")
	repo.Init()

	return dir, repo.(*repositoryGit), shutdown
}

func Test_AnUnitialisedGitRepositoryCanBeRecognisedByAGitInstance(t *testing.T) {

	tempDir, shutdown := tempDirectory(t)
	defer shutdown()

	repo := NewGit(tempDir, "unittest")
	require.False(t, repo.IsInitialised())
}

func Test_AnInitialisedGitRepositoryCanBeRecognisedByAGitInstance(t *testing.T) {

	_, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	require.True(t, repo.IsInitialised())
}

func Test_AnInitialisedGitRepositoryCanHashObjects(t *testing.T) {

	_, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	for _, test := range []struct {
		input  string
		output string
	}{
		{"test", "30d74d258442c7c65512eafab474568dd706c430"},
		{"test2", "d606037cb232bfda7788a8322492312d55b2ae9d"},
		{"some other long string", "5370464603c6098cb422c98b0f3e9a0fdb9c83f8"},
	} {
		t.Run(fmt.Sprintf("input '%s'", test.input), func(t *testing.T) {
			output, err := repo.objectID(test.input)
			require.Nil(t, err)
			require.Equal(t, test.output, output)
		})
	}
}

func Test_AnInitialisedGitRepositoryCanSetMetadata(t *testing.T) {

	_, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	key, data := time.Now().String(), time.Now().String()

	require.Nil(t, repo.SetMetadata(key, data))

	output, err := repo.GetMetadata(key)
	require.Nil(t, err)

	require.Equal(t, data, output)
}
