package cmd

import (
	"fmt"
	"io"

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

func NewCobraHarness(app specstack.SpecStack, stdin io.Reader, stdout, stderr io.Writer) *CobraHarness {
	return &CobraHarness{
		app:    app,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

type CobraHarness struct {
	app    specstack.SpecStack
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (c *CobraHarness) error(returnCode int, err error) error {
	fmt.Fprintf(c.stderr, "Error: %s\n\n", err)

	return NewCliErr(returnCode, err)
}

func (c *CobraHarness) output(msg string) error {
	fmt.Fprintln(c.stdout, msg)

	return nil
}

func (c *CobraHarness) PersistentPreRunE(cmd *cobra.Command, args []string) error {
	if err := c.app.Initialise(); err != nil {
		return c.error(1, err)
	}

	return nil
}

func (c *CobraHarness) ConfigList(cmd *cobra.Command, args []string) error {
	configMap, err := c.app.ListConfiguration()
	if err != nil {
		return c.error(1, err)
	}

	for key, value := range configMap {
		fmt.Fprintf(c.stdout, "%s=%s\n", key, value)
	}

	return nil
}
