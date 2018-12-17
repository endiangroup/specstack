package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/endiangroup/gitest"
	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/config"
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
	repo  *repository.Git
	path  string
	cobra *cobra.Command

	stdout *bytes.Buffer
	stdin  *bytes.Buffer
	stderr *bytes.Buffer

	gitServer *gitest.Server

	assertError error
	exitCode    int
}

func (t *testHarness) ScenarioCleanup(_ interface{}, _ error) {
	if err := t.fs.RemoveAll(t.path); err != nil {
		panic(err)
	}

	if t.gitServer != nil {
		t.gitServer.Close()
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

	if err := t.fs.MkdirAll(filepath.Join(t.path, "features"), 0755); err != nil {
		return err
	}

	return afero.WriteFile(
		t.fs,
		"features/story1.feature",
		[]byte(`Feature: story1`),
		os.ModePerm,
	)
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

func (t *testHarness) iShouldSeeAWarningMessageInformingMe(msg string) error {
	if !assert.Contains(t, t.stderr.String(), msg) {
		t.Errorf("%d\nstdout=%s\nstderr=%s", t.exitCode, t.stdout.String(), t.stderr.String())
		return t.AssertError()
	}

	if !assert.Equal(t, 0, t.exitCode, "Non-zero exit coded returned") {
		return t.AssertError()
	}
	return nil
}

func (t *testHarness) iShouldSeeAHelpfulSuggestionInformingMe(msg string) error {
	if !assert.Contains(t, t.stderr.String(), msg) {
		t.Errorf("%d, %s, %s", t.exitCode, t.stdout.String(), t.stderr.String())
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
	if !assert.True(t, len(strings.Split(t.stdout.String(), "\n")) > 0, "Nothing outputted, expected some lines") {
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

func (t *testHarness) iHaveAStoryCalled(name string) error {
	return t.iHaveAFileCalledWithTheFollowingContent(
		fmt.Sprintf("features/%s.feature", name),
		&gherkin.DocString{
			Content: fmt.Sprintf(`Feature: %s`, name),
		},
	)
}

func (t *testHarness) iHaveAStoryCalledInMySpecWithTheFollowingMetadata(name string, table *gherkin.DataTable) error {
	if err := t.iHaveAStoryCalled(name); err != nil {
		return err
	}
	return t.myStoryHasTheFollowingMetadata(name, table)
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

	return t.AssertMetadataFromStdout(metadataKey, value)
}

func (t *testHarness) theMetadataShouldBeAddedToScenarioWithTheValue(key, scenarioQuery, value string) error {
	if err := t.iRunTheCommand(fmt.Sprintf("metadata ls --scenario %s", scenarioQuery)); err != nil {
		return err
	}

	return t.AssertMetadataFromStdout(key, value)
}

func (t *testHarness) iShouldSeeNoErrors() error {
	if !assert.True(t, t.exitCode == 0, "Non-zero exit coded returned, expected 0") {
		fmt.Println("Stdout:", t.stdout)
		fmt.Println("Stderr:", t.stderr)
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) iHaveAGitinitialisedProjectDirectory() error {
	return t.RunNoArgSteps(
		t.iHaveAProjectDirectory,
		t.iHaveInitialisedGit,
		t.iHaveSetMyUserDetails,
	)
}

func (t *testHarness) iHaveNotConfiguredAProjectRemote() error {
	err := t.iRunTheCommand(`config set project.remote=`)
	return err
}

func (t *testHarness) iHaveNotSetAGitRemote() error {
	t.gitServer = nil
	return nil
}

func (t *testHarness) overwriteHooks() error {
	goPath := os.Getenv("GOPATH")
	cmd := "go run " + filepath.Join(
		goPath,
		"src/github.com/endiangroup/specstack/cmd/spec/*.go",
	)

	if err := t.repo.WriteHookFile("pre-push", cmd+" git-hook exec pre-push"); err != nil {
		return err
	}

	if err := t.repo.WriteHookFile("post-merge", cmd+" git-hook exec post-merge"); err != nil {
		return err
	}

	return nil
}

func (t *testHarness) myStoryHasTheFollowingMetadata(storyName string, table *gherkin.DataTable) error {
	if err := t.thePushingModeIsNotSetToAutomatic(); err != nil {
		return err
	}
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

func (t *testHarness) myScenarioHasTheFollowingMetadata(name string, table *gherkin.DataTable) error {
	if err := t.thePushingModeIsNotSetToAutomatic(); err != nil {
		return err
	}
	for _, row := range table.Rows[1:] {
		if err := t.iRunTheCommand(
			fmt.Sprintf(
				`metadata add --scenario "%s" "%s"="%s"`,
				name,
				row.Cells[0].Value, row.Cells[1].Value,
			),
		); err != nil {
			return err
		}
	}

	return nil
}

func (t *testHarness) iHaveSetThePullingModeToSemiautomatic() error {
	return t.SetSyncMode("pulling", config.ModeSemiAuto)
}

func (t *testHarness) iHaveSetThePullingModeToAutomatic() error {
	return t.SetSyncMode("pulling", config.ModeAuto)
}

func (t *testHarness) iHaveSetThePushingModeToSemiautomatic() error {
	return t.SetSyncMode("pushing", config.ModeSemiAuto)
}

func (t *testHarness) iHaveSetThePushingModeToAutomatic() error {
	return t.SetSyncMode("pushing", config.ModeAuto)
}

func (t *testHarness) thePushingModeIsNotSetToAutomatic() error {
	return t.SetSyncMode("pushing", config.ModeSemiAuto)
}

func (t *testHarness) iAddSomeMetadata() error {
	return t.iRunTheCommand(`metadata add --story story1 key1=value1`)
}

func (t *testHarness) iRunAGitPull() error {
	return t.RunGitCommand("pull")
}

func (t *testHarness) iRunAGitPush() error {
	return t.RunGitCommand("push")
}

func (t *testHarness) iMakeACommit() error {
	if err := afero.WriteFile(
		t.fs,
		"features/story1.feature",
		[]byte(`Feature: Story1 (modified)`),
		os.ModePerm,
	); err != nil {
		return err
	}

	if err := t.RunGitCommand("add", "features/story1.feature"); err != nil {
		return err
	}

	if err := t.RunGitCommand("commit", "-m", "iMakeACommit"); err != nil {
		return err
	}

	return nil
}

func (t *testHarness) iHaveAProperlyConfiguredProjectDirectory() error {
	if err := t.iHaveAnEmptyDirectory(); err != nil {
		return err
	}

	_, f, _, _ := runtime.Caller(1)
	var err error
	t.gitServer, err = gitest.NewServer(filepath.Join(path.Dir(f), "fixtures/git/starting"))
	if err != nil {
		return err
	}

	if err := t.RunGitCommand(
		"clone",
		fmt.Sprintf("%s/%s.git", t.gitServer.URL, t.gitServer.ValidRepo),
		".",
	); err != nil {
		return err
	}

	if err := t.iHaveSetTheGitUserNameTo("Speck Stack"); err != nil {
		return err
	}

	if err := t.iHaveSetTheGitUserEmailTo("dev@specstack.io"); err != nil {
		return err
	}

	return nil
}

func (t *testHarness) theRemoteGitServerIsntRespondingProperly() error {
	t.gitServer.Server.Close()
	return nil
}

func (t *testHarness) iShouldSeeAnAppropriateErrorFromGit() error {
	if !assert.True(t, t.exitCode > 0, "Zero exit coded returned, expected > 0") {
		return t.AssertError()
	}

	if assert.Empty(t, t.stderr.String(), "Expected an error string") {
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) iShouldSeeAnAppropriateWarningFromGit() error {
	if !assert.Equal(t, t.exitCode, 0, "non-zero exit coded returned, expected 0") {
		return t.AssertError()
	}

	if assert.Empty(t, t.stderr.String(), "Expected an error string") {
		return t.AssertError()
	}

	return nil
}

func (t *testHarness) thereAreNewMetadataOnTheRemoteGitServer() error {
	_, f, _, _ := runtime.Caller(1)
	p := filepath.Join(path.Dir(f), "fixtures/git/with-commits")
	return t.gitServer.SetTemplate(p)
}

func (t *testHarness) myMetadataShouldBeFetchedFromTheRemoteGitServer() error {
	if err := t.iRunTheCommand("metadata list --story story1"); err != nil {
		return nil
	}

	scanner := metadata.NewPlaintextPrintscanner()
	entries, err := scanner.Scan(t.stdout)
	if err != nil {
		return nil
	}

	expectedEntries := []metadata.Entry{
		{Name: "a", Value: "a"},
	}

	if !assert.Equal(t, expectedEntries, entries) {
		return fmt.Errorf("Entries not not match as expected")
	}

	return nil
}

func (t *testHarness) myMetadataShouldBePushedToTheRemoteGitServer() error {
	timeout := time.After(10 * time.Millisecond)
	for {
		select {
		case res := <-t.gitServer.RefsEventChan:
			// This is an imperfect test because the mock server
			// doesn't do very much, but we know the client must
			// send git-upload-pack messages to update the remote,
			// so we count that as valid.
			if res.FormValue("service") == "git-upload-pack" {
				return nil
			}
		case <-timeout:
			break
		}
	}
	return fmt.Errorf("Timed out")
}

func (t *testHarness) myStoryHasAScenarioCalledWithTheFollowingMetadata(story, scenario string, table *gherkin.DataTable) error {
	if err := t.iHaveAFileCalledWithTheFollowingContent(
		fmt.Sprintf("features/%s.feature", story),
		&gherkin.DocString{
			Content: fmt.Sprintf(
				`
Feature: %s
    Scenario: %s
	    Then something happens
				`,
				story,
				scenario,
			),
		},
	); err != nil {
		return err
	}

	return t.myScenarioHasTheFollowingMetadata(scenario, table)
}

func (t *testHarness) myStoryHasAScenarioCalledWithSomeMetadata(arg1, arg2 string) error {
	return godog.ErrPending
}

func (t *testHarness) iMakeMinorChangesToScenario(arg1 string) error {
	return godog.ErrPending
}

func (t *testHarness) iCommitAndPushMyChangesWithGit() error {
	return godog.ErrPending
}

func (t *testHarness) theMetadataOnShouldStillExist(arg1 string) error {
	return godog.ErrPending
}

func (t *testHarness) runAnySpecCommand() error {
	return godog.ErrPending
}

func (t *testHarness) Errorf(format string, args ...interface{}) {
	t.assertError = fmt.Errorf(format, args...)
}

func (t *testHarness) AssertError() error {
	return t.assertError
}

func (t *testHarness) RunNoArgSteps(fns ...func() error) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

func (t *testHarness) RunGitCommand(args ...string) error {

	cmd := exec.Command("git", args...)
	cmd.Dir = t.path
	cmd.Stdout = t.stdout
	cmd.Stderr = t.stderr
	cmd.Stdin = t.stdin

	err := cmd.Run()

	if err != nil {
		t.exitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				t.exitCode = status.ExitStatus()
			}
		}
	}

	return err
}

func (t *testHarness) SetSyncMode(mode, value string) error {
	if err := t.overwriteHooks(); err != nil {
		return err
	}
	return t.iRunTheCommand(fmt.Sprintf(`config set project.%smode=%s`, mode, value))
}

func (t *testHarness) AssertMetadataFromStdout(key, value string) error {
	scanner := metadata.NewPlaintextPrintscanner()
	entries, err := scanner.Scan(t.stdout)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.Name == key {
			if entry.Value == value {
				return nil
			} else {
				return fmt.Errorf("Got %s, expected %s", entry.Value, value)
			}
		}
	}

	return fmt.Errorf("metadata not found")
}

func FeatureContext(s *godog.Suite) {
	th := newTestHarness()

	s.Step(`^I have an empty directory$`, th.iHaveAnEmptyDirectory)
	s.Step(`^I have a project directory$`, th.iHaveAProjectDirectory)
	s.Step(`^I run "([^"]*)"$`, th.iRunTheCommand)
	s.Step(`^I should see an error message informing me "([^"]*)"$`, th.iShouldSeeAnErrorMessageInformingMe)
	s.Step(`^I should see a warning message informing me "([^"]*)"$`, th.iShouldSeeAWarningMessageInformingMe)
	s.Step(`^I should see a helpful suggestion informing me "([^"]*)"$`, th.iShouldSeeAHelpfulSuggestionInformingMe)
	s.Step(`^I have initialised git$`, th.iHaveInitialisedGit)
	s.Step(`^I should see the following:$`, th.iShouldSeeTheFollowing)
	s.Step(`^I should see some configuration keys and values$`, th.iShouldSeeSomeConfigurationKeysAndValues)
	s.Step(`^The config key "([^"]*)" should equal "([^"]*)"$`, th.theConfigKeyShouldEqual)
	s.Step(`^I have no user details$`, th.iHaveNoUserDetails)
	s.Step(`^I have set the git user name to "([^"]*)"$`, th.iHaveSetTheGitUserNameTo)
	s.Step(`^I have set the git user email to "([^"]*)"$`, th.iHaveSetTheGitUserEmailTo)
	s.Step(`^I have set my user details$`, th.iHaveSetMyUserDetails)
	s.Step(`^I have a file called "([^"]*)" with the following content:$`, th.iHaveAFileCalledWithTheFollowingContent)
	s.Step(`^I have a story called "([^"]*)"$`, th.iHaveAStoryCalled)
	s.Step(`^I have a story called "([^"]*)" in my spec with the following metadata:$`, th.iHaveAStoryCalledInMySpecWithTheFollowingMetadata)
	s.Step(`^My story "([^"]*)" has a scenario called "([^"]*)" with the following metadata:$`, th.myStoryHasAScenarioCalledWithTheFollowingMetadata)
	s.Step(`^My story "([^"]*)" has a scenario called "([^"]*)" with some metadata$`, th.myStoryHasAScenarioCalledWithSomeMetadata)
	s.Step(`^My story "([^"]*)" has the following metadata:$`, th.myStoryHasTheFollowingMetadata)
	s.Step(`^I have configured git$`, th.iHaveConfiguredGit)
	s.Step(`^I have not initialised git$`, th.iHaveNotInitialisedGit)
	s.Step(`^I have a configured project directory$`, th.iHaveAConfiguredProjectDirectory)
	s.Step(`^The metadata "([^"]*)" should be added to story "([^"]*)" with the value "([^"]*)"$`, th.theMetadataShouldBeAddedToStory)
	s.Step(`^The metadata "([^"]*)" should be added to scenario "([^"]*)" with the value "([^"]*)"$`, th.theMetadataShouldBeAddedToScenarioWithTheValue)
	s.Step(`^I should see no errors$`, th.iShouldSeeNoErrors)
	s.Step(`^I have a git-initialised project directory$`, th.iHaveAGitinitialisedProjectDirectory)
	s.Step(`^I have not configured a project remote$`, th.iHaveNotConfiguredAProjectRemote)
	s.Step(`^I have not set a git remote$`, th.iHaveNotSetAGitRemote)
	s.Step(`^I have set the pulling mode to semi-automatic$`, th.iHaveSetThePullingModeToSemiautomatic)
	s.Step(`^I add some metadata$`, th.iAddSomeMetadata)
	s.Step(`^I have added some metadata$`, th.iAddSomeMetadata)
	s.Step(`^I run a git pull$`, th.iRunAGitPull)
	s.Step(`^I run a git push$`, th.iRunAGitPush)
	s.Step(`^I make a commit$`, th.iMakeACommit)
	s.Step(`^I have set the pulling mode to automatic$`, th.iHaveSetThePullingModeToAutomatic)
	s.Step(`^I have set the pushing mode to semi-automatic$`, th.iHaveSetThePushingModeToSemiautomatic)
	s.Step(`^I have set the pushing mode to automatic$`, th.iHaveSetThePushingModeToAutomatic)
	s.Step(`^the pushing mode is not set to automatic$`, th.thePushingModeIsNotSetToAutomatic)
	s.Step(`^I have a properly configured project directory$`, th.iHaveAProperlyConfiguredProjectDirectory)
	s.Step(`^The remote git server isn\'t responding properly$`, th.theRemoteGitServerIsntRespondingProperly)
	s.Step(`^I should see an appropriate error from git$`, th.iShouldSeeAnAppropriateErrorFromGit)
	s.Step(`^I should see an appropriate warning from git$`, th.iShouldSeeAnAppropriateWarningFromGit)
	s.Step(`^there are new metadata on the remote git server$`, th.thereAreNewMetadataOnTheRemoteGitServer)
	s.Step(`^my metadata should be fetched from the remote git server$`, th.myMetadataShouldBeFetchedFromTheRemoteGitServer)
	s.Step(`^my metadata should be pushed to the remote git server$`, th.myMetadataShouldBePushedToTheRemoteGitServer)
	s.Step(`^I make minor changes to scenario "([^"]*)"$`, th.iMakeMinorChangesToScenario)
	s.Step(`^I commit and push my changes with git$`, th.iCommitAndPushMyChangesWithGit)
	s.Step(`^the metadata on "([^"]*)" should still exist$`, th.theMetadataOnShouldStillExist)
	s.Step(`^I run any spec command$`, th.runAnySpecCommand)

	s.AfterScenario(th.ScenarioCleanup)
}
