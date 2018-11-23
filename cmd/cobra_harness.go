package cmd

import (
	"fmt"
	"io"
	"strings"

	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/metadata"
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

func (c *CobraHarness) errorOrNil(cmd *cobra.Command, returnCode int, err error) error {
	if err == nil {
		return nil
	}

	cmd.Root().SetOutput(c.stderr)

	return NewCliErr(returnCode, err)
}

func (c *CobraHarness) flagValueString(cmd *cobra.Command, name string) string {
	return cmd.Flag(name).Value.String()
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

func (c *CobraHarness) SetKeyValueArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
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
	entityFound := false

	if storyName := c.flagValueString(cmd, "story"); storyName != "" {
		entityFound = true
		for _, arg := range args {
			kv := strings.Split(arg, "=")
			if err := c.app.AddMetadataToStory(storyName, kv[0], kv[1]); err != nil {
				return c.error(cmd, 1, err)
			}
		}
	}

	if !entityFound {
		return c.error(cmd, 0, fmt.Errorf("specify a story"))
	}

	return nil
}

func (c *CobraHarness) MetadataList(cmd *cobra.Command, args []string) error {
	var entries []*metadata.Entry
	entityFound := false

	if storyName := c.flagValueString(cmd, "story"); storyName != "" {
		var err error
		entries, err = c.app.GetStoryMetadata(storyName)
		if err != nil {
			return c.error(cmd, 1, err)
		}
		entityFound = true
	}

	if !entityFound {
		return c.error(cmd, 0, fmt.Errorf("specify a story"))
	}

	printer := metadata.NewPlaintextPrintscanner()
	return printer.Print(c.stdout, entries)
}

func (c *CobraHarness) GitHookExec(cmd *cobra.Command, args []string) error {
	switch args[0] {
	case "pre-commit":
		return c.errorOrNil(cmd, 1, c.app.RunRepoPreCommitHook())

	case "post-commit":
		return c.errorOrNil(cmd, 1, c.app.RunRepoPostCommitHook())
	}

	return c.error(cmd, 1, fmt.Errorf("invalid hook name"))
}

func (c *CobraHarness) Pull(cmd *cobra.Command, args []string) error {
	return c.errorOrNil(cmd, 1, c.app.Pull())
}

func (c *CobraHarness) Push(cmd *cobra.Command, args []string) error {
	return c.errorOrNil(cmd, 1, c.app.Push())
}
