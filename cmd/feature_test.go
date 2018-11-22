package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/metadata"
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

	git := repository.NewGitRepository(tmpPath, repository.GitConfigScopeLocal)
	th.repo = git

	repoStore := persistence.NewStore(
		persistence.NewNamespacedKeyValueStorer(th.repo, "specstack"),
		git,
	)
	developer := personas.NewDeveloper(repoStore)
	app := specstack.New(testdirPath, th.repo, developer, repoStore, th.stdout, th.stderr)

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

func (t *testHarness) iHaveAProjectDirectory() error {
	if err := t.iHaveAnEmptyDirectory(); err != nil {
		return nil
	}

	return t.fs.MkdirAll(filepath.Join(t.path, "features"), 0755)
}

func (t *testHarness) iRunTheCommand(cmd string) error {
	r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)

	processed := []string{}
	concat := false
	index := 0
	for i, arg := range r.FindAllString(cmd, -1) {
		value := strings.Trim(arg, `"`)
		if i > 0 && (arg == "=" || concat) {
			processed[index-1] += value
			concat = true
			continue
		}
		processed = append(processed, value)
		concat = false
		index++
	}

	t.cobra.SetArgs(processed)
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
		t.Errorf("%d, %s, %s", t.exitCode, t.stdout.String(), t.stderr.String())
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

func (t *testHarness) iHaveNotInitialisedGit() error {
	return nil
}

func (t *testHarness) iHaveConfiguredGit() error {

	if err := t.iHaveInitialisedGit(); err != nil {
		return nil
	}

	if err := t.iHaveSetTheGitUserNameTo("Speck Stack"); err != nil {
		return err
	}

	if err := t.iHaveSetTheGitUserEmailTo("dev@specstack.io"); err != nil {
		return err
	}

	return nil
}

func (t *testHarness) iShouldSeeTheFollowing(output *gherkin.DocString) error {
	lines := strings.Split(output.Content, "\n")
	for _, line := range lines {
		if !assert.Contains(t, t.stdout.String(), line) {
			return t.AssertError()
		}
	}

	return t.iShouldSeeNoErrors()
}

func (t *testHarness) iShouldSeeSomeConfigurationKeysAndValues() error {
	if !assert.True(t, len(strings.Split(t.stdout.String(), "\n")) > 0, "Nothing outputed, expected some lines") {
		return t.AssertError()
	}

	if !assert.Regexp(t, `[a-z.]+=.+(\n)?`, t.stdout) {
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) theConfigKeyShouldEqual(key, expectedValue string) error {
	value, err := t.repo.GetConfig("specstack." + key)
	if err != nil {
		return err
	}

	if !assert.Equal(t, expectedValue, value) {
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) iHaveNoUserDetails() error {
	if err := t.repo.UnsetConfig("user.name"); err != nil {
		return err
	}
	if err := t.repo.UnsetConfig("user.email"); err != nil {
		return err
	}

	return nil
}

func (t *testHarness) iHaveSetTheGitUserNameTo(name string) error {
	return t.repo.SetConfig("user.name", name)
}

func (t *testHarness) iHaveSetTheGitUserEmailTo(email string) error {
	return t.repo.SetConfig("user.email", email)
}

func (t *testHarness) iHaveSetMyUserDetails() error {
	err := t.iHaveSetTheGitUserNameTo("Spec Stack")
	if err != nil {
		return err
	}

	return t.iHaveSetTheGitUserEmailTo("dev@specstack.io")
}

func (t *testHarness) iHaveAFileCalledWithTheFollowingContent(filename string, content *gherkin.DocString) error {
	if err := t.fs.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return afero.WriteFile(t.fs, filename, []byte(content.Content), os.ModePerm)
}

func (t *testHarness) iHaveAConfiguredProjectDirectory() error {
	if err := t.iHaveAProjectDirectory(); err != nil {
		return err
	}
	return t.iHaveConfiguredGit()
}

func (t *testHarness) theMetadataShouldBeAddedToStory(metadataKey, storyId, value string) error {
	if err := t.iRunTheCommand(fmt.Sprintf("metadata ls --story %s", storyId)); err != nil {
		return err
	}

	scanner := metadata.NewPlaintextPrintscanner()
	entries, err := scanner.Scan(t.stdout)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name == metadataKey {
			if entry.Value == value {
				return nil
			} else {
				return fmt.Errorf("Got %s, expected %s", entry.Value, value)
			}
		}
	}

	return fmt.Errorf("metadata not found")
}

func (t *testHarness) iShouldSeeNoErrors() error {
	if !assert.True(t, t.exitCode == 0, "Non-zero exit coded returned, expected 0") {
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) hasTheFollowingMetadata(storyName string, table *gherkin.DataTable) error {
	for _, row := range table.Rows[1:] {
		if err := t.iRunTheCommand(
			fmt.Sprintf(
				`metadata add --story "%s" "%s"="%s"`,
				storyName,
				row.Cells[0].Value, row.Cells[1].Value,
			),
		); err != nil {
			return err
		}
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
	s.Step(`^I have a project directory$`, th.iHaveAProjectDirectory)
	s.Step(`^I run "([^"]*)"$`, th.iRunTheCommand)
	s.Step(`^I should see an error message informing me "([^"]*)"$`, th.iShouldSeeAnErrorMessageInformingMe)
	s.Step(`^I have initialised git$`, th.iHaveInitialisedGit)
	s.Step(`^I should see the following:$`, th.iShouldSeeTheFollowing)
	s.Step(`^I should see some configuration keys and values$`, th.iShouldSeeSomeConfigurationKeysAndValues)
	s.Step(`^The config key "([^"]*)" should equal "([^"]*)"$`, th.theConfigKeyShouldEqual)
	s.Step(`^I have no user details$`, th.iHaveNoUserDetails)
	s.Step(`^I have set the git user name to "([^"]*)"$`, th.iHaveSetTheGitUserNameTo)
	s.Step(`^I have set the git user email to "([^"]*)"$`, th.iHaveSetTheGitUserEmailTo)
	s.Step(`^I have set my user details$`, th.iHaveSetMyUserDetails)
	s.Step(`^I have a file called "([^"]*)" with the following content:$`, th.iHaveAFileCalledWithTheFollowingContent)
	s.Step(`^I have configured git$`, th.iHaveConfiguredGit)
	s.Step(`^I have not initialised git$`, th.iHaveNotInitialisedGit)
	s.Step(`^I have a configured project directory$`, th.iHaveAConfiguredProjectDirectory)
	s.Step(`^The metadata "([^"]*)" should be added to story "([^"]*)" with the value "([^"]*)"$`, th.theMetadataShouldBeAddedToStory)
	s.Step(`^I should see no errors$`, th.iShouldSeeNoErrors)
	s.Step(`^"([^"]*)" has the following metadata:$`, th.hasTheFollowingMetadata)

	s.AfterScenario(th.ScenarioCleanup)
}
