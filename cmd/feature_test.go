package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/actors"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/repository"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func newTestHarness() *testHarness {
	// Current git implementation can only interact with OS FS
	fs := afero.NewOsFs()

	tmpPath, err := afero.TempDir(fs, "", "specstack-")
	if err != nil {
		panic(err)
	}
	testdirPath := filepath.Join(tmpPath, "test-dir")

	th := &testHarness{
		fs:     fs,
		path:   tmpPath,
		stdout: bytes.NewBuffer(nil),
		stdin:  bytes.NewBuffer(nil),
		stderr: bytes.NewBuffer(nil),
	}

	th.repo = repository.NewGitRepo(tmpPath, "specstack")
	repoStore := persistence.NewRepositoryStore(th.repo)
	developer := actors.NewDeveloper(repoStore)
	app := specstack.NewApp(testdirPath, th.repo, developer, repoStore)

	th.cobra = WireUpHarness(NewCobraHarness(app, th.stdin, th.stdout, th.stderr))

	return th
}

type testHarness struct {
	fs    afero.Fs
	repo  repository.ReadWriter
	path  string
	cobra *cobra.Command

	stdout *bytes.Buffer
	stdin  *bytes.Buffer
	stderr *bytes.Buffer

	assertError error
	exitCode    int
}

func (t *testHarness) ScenarioCleanup(_ interface{}, _ error) {
	if err := t.fs.RemoveAll(t.path); err != nil {
		panic(err)
	}

	*t = *newTestHarness()
}

func (t *testHarness) iHaveAnEmptyDirectory() error {
	if err := t.fs.MkdirAll(t.path, 0755); err != nil {
		return err
	}

	return os.Chdir(t.path)
}

func (t *testHarness) iRunTheCommand(cmd string) error {
	t.cobra.SetArgs(strings.Split(cmd, " "))
	err := t.cobra.Execute()
	if err != nil {
		if cliErr, ok := err.(CliErr); ok {
			t.exitCode = cliErr.ExitCode
		}
	}

	return nil
}

func (t *testHarness) iShouldSeeAnErrorMessageInformingMe(msg string) error {
	if !assert.Contains(t, t.stderr.String(), msg) {
		return t.AssertError()
	}

	if !assert.True(t, t.exitCode > 0, "Zero exit coded returned, expected > 0") {
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) iHaveInitialisedGit() error {
	return t.repo.Init()
}

func (t *testHarness) iShouldSeeTheFollowing(output *gherkin.DocString) error {
	lines := strings.Split(output.Content, "\n")
	for _, line := range lines {
		if !assert.Contains(t, t.stdout.String(), line) {
			return t.AssertError()
		}
	}

	if !assert.True(t, t.exitCode == 0, "Non-zero exit coded returned, expected 0") {
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) Errorf(format string, args ...interface{}) {
	t.assertError = fmt.Errorf(format, args...)
}

func (t *testHarness) AssertError() error {
	return t.assertError
}

func FeatureContext(s *godog.Suite) {
	th := newTestHarness()

	s.Step(`^I have an empty directory$`, th.iHaveAnEmptyDirectory)
	s.Step(`^I run "([^"]*)"$`, th.iRunTheCommand)
	s.Step(`^I should see an error message informing me "([^"]*)"$`, th.iShouldSeeAnErrorMessageInformingMe)
	s.Step(`^I have initialised git$`, th.iHaveInitialisedGit)
	s.Step(`^I should see the following:$`, th.iShouldSeeTheFollowing)

	s.AfterScenario(th.ScenarioCleanup)
}
