package specstack

import (
	"bytes"
	"fmt"

	"github.com/DATA-DOG/godog"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func newTestHarness() *testHarness {
	th := &testHarness{
		fs: afero.NewMemMapFs(),
	}

	th.app = NewCliApp(th.fs)

	return th
}

type testHarness struct {
	app CliApp
	fs  afero.Fs

	stdout bytes.Buffer
	stdin  bytes.Buffer
	stderr bytes.Buffer

	assertError error
	returnCode  int
}

func (t *testHarness) ScenarioCleanup(_ interface{}, _ error) {
	*t = *newTestHarness()
}

func (t *testHarness) iHaveAnEmptyDirectory() error {
	return t.fs.Mkdir("test-dir", 0755)
}

func (t *testHarness) iRunTheCommand(cmd string) error {
	t.returnCode = t.app.Run([]string{"init"}, &t.stdout, &t.stdin, &t.stderr)

	return nil
}

func (t *testHarness) iShouldSeeAnErrorMessageInformingMe(msg string) error {
	if !assert.Contains(t, t.stderr.String(), msg) {
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
	s.Step(`^I run the "([^"]*)" command$`, th.iRunTheCommand)
	s.Step(`^I should see an error message informing me "([^"]*)"$`, th.iShouldSeeAnErrorMessageInformingMe)

	s.AfterScenario(th.ScenarioCleanup)
}
