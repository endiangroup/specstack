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
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/personas"
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

	th.repo = repository.NewGit(tmpPath, "specstack")
	repoStore := persistence.NewRepositoryStore(th.repo)
	developer := personas.NewDeveloper(repoStore)
	app := specstack.New(testdirPath, th.repo, developer, repoStore)

	th.cobra = WireUpCobraHarness(NewCobraHarness(app, th.stdin, th.stdout, th.stderr))

	return th
}

type testHarness struct {
	fs    afero.Fs
	repo  repository.Repository
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

func (t *testHarness) iShouldSeeSomeConfigurationKeysAndValues() error {
	if !assert.True(t, len(strings.Split(t.stdout.String(), "\n")) > 0, "Nothing outputed, expected some lines") {
		return t.AssertError()
	}

	if !assert.Regexp(t, `[a-z.]+=.+(\n)?`, t.stdout) {
		t.AssertError()
	}

	return nil
}

func (t *testHarness) theConfigKeyShouldEqual(key, expectedValue string) error {
	value, err := t.repo.Get(key)
	if err != nil {
		return err
	}

	if !assert.Equal(t, expectedValue, value) {
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
	s.Step(`^I should see some configuration keys and values$`, th.iShouldSeeSomeConfigurationKeysAndValues)
	s.Step(`^The config key "([^"]*)" should equal "([^"]*)"$`, th.theConfigKeyShouldEqual)

	s.AfterScenario(th.ScenarioCleanup)
}
