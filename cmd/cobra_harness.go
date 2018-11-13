package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/endiangroup/specstack"
	"github.com/spf13/cobra"
)

func NewCliErr(exitCode int, err error) CliErr {
	return CliErr{ExitCode: exitCode, Err: err}
}

type CliErr struct {
	ExitCode int
	Err      error
}

func (err CliErr) Error() string {
	return err.Err.Error()
}

func NewCobraHarness(app specstack.Controller, stdin io.Reader, stdout, stderr io.Writer) *CobraHarness {
	return &CobraHarness{
		app:    app,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

type CobraHarness struct {
	app    specstack.Controller
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *CobraHarness) error(cmd *cobra.Command, returnCode int, err error) error {
	cmd.Root().SetOutput(c.stderr)

	return NewCliErr(returnCode, err)
}

func (c *CobraHarness) warning(message string) {
	c.stdout.Write([]byte(fmt.Sprintf("WARNING: %s\n", message)))
}

func (c *CobraHarness) PersistentPreRunE(cmd *cobra.Command, args []string) error {
	if err := c.app.Initialise(); err != nil {
		return c.error(cmd, 1, err)
	}

	return nil
}

func (c *CobraHarness) ConfigList(cmd *cobra.Command, args []string) error {
	configMap, err := c.app.ListConfiguration()
	if err != nil {
		return c.error(cmd, 1, err)
	}

	for key, value := range configMap {
		cmd.Printf("%s=%s\n", key, value)
	}

	return nil
}

func (c *CobraHarness) ConfigGet(cmd *cobra.Command, args []string) error {
	value, err := c.app.GetConfiguration(args[0])
	if err != nil {
		return c.error(cmd, 1, err)
	}

	cmd.Print(value)

	return nil
}

func (c *CobraHarness) ConfigSetArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.ExactArgs(1)(cmd, args); err != nil {
		return c.error(cmd, 1, err)
	}

	if err := IsKeyEqualsValueFormat(args[0]); err != nil {
		return c.error(cmd, 1, err)
	}

	return nil

}

func (c *CobraHarness) ConfigSet(cmd *cobra.Command, args []string) error {
	keyValueParts := strings.Split(args[0], "=")

	err := c.app.SetConfiguration(keyValueParts[0], keyValueParts[1])
	if err != nil {
		return c.error(cmd, 1, err)
	}

	return nil
}

func (c *CobraHarness) MetadataAdd(cmd *cobra.Command, args []string) error {

	spec, warnings, err := c.app.Specification()

	if err != nil {
		return c.error(cmd, 1, err)
	}

	for _, warning := range warnings {
		c.warning(warning.Error())
	}

	storyFlag := cmd.Flag("story")

	if storyName := storyFlag.Value.String(); storyName != "" {
		story, err := spec.FindStory(storyName)

		if err != nil {
			return c.error(cmd, 1, err)
		}

		return fmt.Errorf("TODO: add metadata: %s", story.Name)
	}

	return nil
}
