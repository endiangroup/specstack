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
	repo := NewGitRepository(dir)

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

func getFileMetadata(t *testing.T, repo *Git, fileName string) (value []string) {
	f0, err := os.Open(fileName)
	require.NotNil(t, f0)
	require.Nil(t, err)
	require.Nil(t, repo.GetMetadata(f0, &value))
	require.Nil(t, f0.Close())
	return
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

	t.Run("Get correct output", func(t *testing.T) {
		output := []string{}
		require.Nil(t, repo.GetMetadata(bytes.NewBufferString(key), &output))
		require.Equal(t, []string{data}, output)
	})

	t.Run("Get nothing when there's no note", func(t *testing.T) {
		output2 := []string{}
		require.Nil(t, repo.GetMetadata(bytes.NewBufferString("doesn't exist"), &output2))
		require.Equal(t, []string{}, output2)
	})
}

func Test_AnInitialisedGitRepositoryCanSetComplexMetadata(t *testing.T) {

	_, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	type myStruct struct {
		A int
		B int
	}

	key, data := time.Now().String(), myStruct{}

	require.Nil(t, repo.SetMetadata(bytes.NewBufferString(key), data))

	t.Run("Get correct output", func(t *testing.T) {
		output := []myStruct{}
		require.Nil(t, repo.GetMetadata(bytes.NewBufferString(key), &output))
		require.Equal(t, []myStruct{data}, output)
	})

	t.Run("Get nothing when there's no note", func(t *testing.T) {
		output := []myStruct{}
		require.Nil(t, repo.GetMetadata(bytes.NewBufferString("doesn't exist"), &output))
		require.Equal(t, []myStruct{}, output)
	})

	t.Run("Get error when there's a type mismatch", func(t *testing.T) {
		output := []string{}
		err := repo.GetMetadata(bytes.NewBufferString(key), &output)
		require.NotNil(t, err)
		require.Equal(t, "json: cannot unmarshal object into Go value of type string", err.Error())
	})
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
		assertGitCmd(t, repo, "", "add", "a.txt")
		assertGitCmd(t, repo, "", "add", "b.txt")
		assertGitCmd(t, repo, "", "commit", "-m", "Commit c")
	})

	t.Run("Check the metadata after commit", func(t *testing.T) {
		require.Equal(t, []string{"m0", "m1"}, getFileMetadata(t, repo, "b.txt"))
	})

	t.Run("Add some more metadata whtout a commit", func(t *testing.T) {
		setFileMetadata(t, repo, "b.txt", "m2")
		setFileMetadata(t, repo, "b.txt", "m3")
		require.Equal(t, []string{"m0", "m1", "m2", "m3"}, getFileMetadata(t, repo, "b.txt"))
	})
}
