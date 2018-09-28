package cmd

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/actors"
	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/repository"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func newTestHarness() *testHarness {
	// Current git implementation can only interact with OS FS
	fs := afero.NewOsFs()

	path, err := afero.TempDir(fs, "", "specstack-")
	if err != nil {
		panic(err)
	}

	th := &testHarness{
		fs:   fs,
		path: path,
	}

	th.repo = repository.NewGitRepo(path)
	config := config.NewRepositoryConfig(th.repo)
	developer := actors.NewDeveloper(config)
	app := specstack.NewApp(th.repo, developer)

	WireUpHarness(NewCobraHarness(app, &th.stdin, &th.stdout, &th.stderr))

	return th
}

type testHarness struct {
	fs   afero.Fs
	repo repository.ReadWriter
	path string

	stdout bytes.Buffer
	stdin  bytes.Buffer
	stderr bytes.Buffer

	assertError error
	returnCode  int
}

func (t *testHarness) ScenarioCleanup(_ interface{}, _ error) {
	if err := t.fs.RemoveAll(t.path); err != nil {
		panic(err)
	}

	*t = *newTestHarness()
}

func (t *testHarness) iHaveAnEmptyDirectory() error {
	return t.fs.MkdirAll("test-dir", 0755)
}

func (t *testHarness) iRunTheCommand(cmd string) error {
	Root.SetArgs(strings.Split(cmd, " "))
	err := Root.Execute()
	if err != nil {
		if cliErr, ok := err.(CliErr); ok {
			t.returnCode = cliErr.ReturnCode
		}
	}

	return nil
}

func (t *testHarness) iShouldSeeAnErrorMessageInformingMe(msg string) error {
	if !assert.Contains(t, t.stderr.String(), msg) {
		return t.AssertError()
	}

	if !assert.True(t, t.returnCode > 0, "Zero return coded returned, expected > 0") {
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) iHaveInitialisedGit() error {
	return t.repo.Init()
}

func (t *testHarness) iShouldSeeTheFollowing(output *gherkin.DocString) error {
	if !assert.Contains(t, t.stdout.String(), output) {
		return t.AssertError()
	}

	if !assert.True(t, t.returnCode > 0, "Zero return coded returned, expected > 0") {
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
