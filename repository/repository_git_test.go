package repository

import (
	"bytes"
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

	require.Nil(t, os.Chdir(dir))

	return dir, func() {
		require.Nil(t, os.RemoveAll(dir))
	}
}

func initialisedGitRepoDir(t *testing.T) (path string, r *Git, shutdown func()) {

	dir, shutdown := tempDirectory(t)
	repo := NewGitRepository(dir).(*Git)

	require.Nil(t, repo.Init())

	_, err := repo.runGitCommand("config", "user.name", "SpecStack")
	require.Nil(t, err)

	_, err = repo.runGitCommand("config", "user.email", "test@specstack.io")
	require.Nil(t, err)

	return dir, repo, shutdown
}

func assertGitCmd(t *testing.T, repo *Git, expectedOutput string, input ...string) {

	output, err := repo.runGitCommand(input...)
	require.Nil(t, err)

	if expectedOutput != "" {
		require.Equal(t, expectedOutput, output)
	}
}

func setFileMetadata(t *testing.T, repo *Git, fileName, value string) {
	f0, err := os.Open(fileName)
	require.NotNil(t, f0)
	require.Nil(t, err)
	require.Nil(t, repo.SetMetadata(f0, value))
	require.Nil(t, f0.Close())
}

func getFileMetadata(t *testing.T, repo *Git, fileName string) []string {
	f0, err := os.Open(fileName)
	require.NotNil(t, f0)
	require.Nil(t, err)
	value, err := repo.GetMetadata(f0)
	require.Nil(t, err)
	require.Nil(t, f0.Close())
	return value
}

func Test_AnUnitialisedGitRepositoryCanBeRecognisedByAGitInstance(t *testing.T) {

	tempDir, shutdown := tempDirectory(t)
	defer shutdown()

	repo := NewGitRepository(tempDir)
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
			output, err := repo.objectID(bytes.NewBufferString(test.input))
			require.Nil(t, err)
			require.Equal(t, test.output, output)
		})
	}
}

func Test_AnInitialisedGitRepositoryCanSetBasicMetadata(t *testing.T) {

	_, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	key, data := time.Now().String(), time.Now().String()

	require.Nil(t, repo.SetMetadata(bytes.NewBufferString(key), data))

	output, err := repo.GetMetadata(bytes.NewBufferString(key))
	require.Nil(t, err)
	require.Equal(t, []string{data}, output)

	output, err = repo.GetMetadata(bytes.NewBufferString("doesn't exist"))
	require.Nil(t, err)
	require.Equal(t, []string{}, output)
}

func Test_AnInitialisedGitRepoCanTrackAndSetMetadataAtTheFileLevel(t *testing.T) {

	_, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	t.Run("Create a file", func(t *testing.T) {
		require.Nil(t, ioutil.WriteFile("a.txt", []byte("1"), os.ModePerm))
	})

	t.Run("Commit the file", func(t *testing.T) {
		assertGitCmd(t, repo, "", "add", "a.txt")
		assertGitCmd(t, repo, "", "commit", "-m", "Commit A")
	})

	t.Run("Set some metadata", func(t *testing.T) {
		setFileMetadata(t, repo, "a.txt", "m0")
	})

	t.Run("Verify the metadata", func(t *testing.T) {
		require.Equal(t, getFileMetadata(t, repo, "a.txt"), []string{"m0"})
	})

	t.Run("Change the file", func(t *testing.T) {
		require.Nil(t, ioutil.WriteFile("a.txt", []byte("2"), os.ModePerm))
	})

	t.Run("Commit the file again", func(t *testing.T) {
		assertGitCmd(t, repo, "", "add", "a.txt")
		assertGitCmd(t, repo, "", "commit", "-m", "Commit B")
	})

	t.Run("Check the metadata after commit", func(t *testing.T) {
		require.Equal(t, []string{"m0"}, getFileMetadata(t, repo, "a.txt"))
	})

	t.Run("Add some more metadata", func(t *testing.T) {
		setFileMetadata(t, repo, "a.txt", "m1")
		require.Equal(t, []string{"m0", "m1"}, getFileMetadata(t, repo, "a.txt"))
	})

	t.Run("Rename the file", func(t *testing.T) {
		require.Nil(t, os.Rename("a.txt", "b.txt"))
		require.Equal(t, []string{"m0", "m1"}, getFileMetadata(t, repo, "b.txt"))
	})

	t.Run("Commit the file again", func(t *testing.T) {
		assertGitCmd(t, repo, "", "add", "b.txt")
		assertGitCmd(t, repo, "", "commit", "-m", "Commit c")
	})

	t.Run("Check the metadata after commit", func(t *testing.T) {
		require.Equal(t, []string{"m0", "m1"}, getFileMetadata(t, repo, "b.txt"))
	})
}
