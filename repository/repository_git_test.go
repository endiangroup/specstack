package repository

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	require.Nil(t, repo.SetMetadata(f0, []byte(value)))
	require.Nil(t, f0.Close())
}

func getFileMetadata(t *testing.T, repo *Git, fileName string) (value []string) {
	f0, err := os.Open(fileName)
	require.NotNil(t, f0)
	require.Nil(t, err)

	raw, err := repo.GetMetadata(f0)
	require.Nil(t, err)

	for _, v := range raw {
		value = append(value, string(v))
	}

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

	key, data := time.Now().String(), []byte(time.Now().String())

	require.Nil(t, repo.SetMetadata(bytes.NewBufferString(key), data))

	t.Run("Get correct output", func(t *testing.T) {
		output, err := repo.GetMetadata(bytes.NewBufferString(key))
		require.Nil(t, err)
		require.Equal(t, [][]byte{data}, output)
	})

	t.Run("Get nothing when there's no note", func(t *testing.T) {
		output, err := repo.GetMetadata(bytes.NewBufferString("doesn't exist"))
		require.Nil(t, err)
		require.Equal(t, [][]byte{}, output)
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

func Test_AnInitialisedGitRepoThrowsAnErrorOnNoRemotePullAndPush(t *testing.T) {

	_, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	t.Run("Pull", func(t *testing.T) {
		err := repo.PullMetadata("doesntexist")
		require.NotNil(t, err)
		require.Equal(t, "fatal: No such remote 'doesntexist'", err.Error())
	})
	t.Run("Push", func(t *testing.T) {
		err := repo.PushMetadata("doesntexist")
		require.NotNil(t, err)
		require.Equal(t, "fatal: No such remote 'doesntexist'", err.Error())
	})
}

func Test_AnInitialisedGitRepoKnowsItsGitDirectories(t *testing.T) {

	dir, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	expectedGitDir := filepath.Join(dir, ".git")
	expectedHooksDir := filepath.Join(expectedGitDir, "hooks")

	topDir, err := repo.topDirectory()
	require.Nil(t, err)
	require.Equal(t, dir, topDir)

	gitDir, err := repo.gitDirectory()
	require.Nil(t, err)
	require.Equal(t, expectedGitDir, gitDir)

	hooksDir, err := repo.gitHooksDirectory()
	require.Nil(t, err)
	require.Equal(t, expectedHooksDir, hooksDir)
}

func Test_AnInitialisedGitRepoCanWriteItsHooksWhenAppropriate(t *testing.T) {

	_, repo, shutdown := initialisedGitRepoDir(t)
	defer shutdown()

	hooksDir, err := repo.gitHooksDirectory()
	require.Nil(t, err)

	pc, pu := filepath.Join(hooksDir, "post-commit"), filepath.Join(hooksDir, "post-update")

	t.Run("Make sure hooks don't exist initially", func(t *testing.T) {
		_, err := os.Stat(pc)
		require.True(t, os.IsNotExist(err))

		_, err = os.Stat(pu)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("Hooks shoulf exist afrer preparation", func(t *testing.T) {
		require.Nil(t, repo.PrepareMetadataSync())

		_, err := os.Stat(pc)
		require.False(t, os.IsNotExist(err))

		_, err = os.Stat(pu)
		require.False(t, os.IsNotExist(err))
	})
}
