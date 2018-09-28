package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/endiangroup/specstack"
	"github.com/spf13/cobra"
)

var (
	ErrUninitialisedRepo = errors.New("Please initialise repository first before running")
)

func NewCliErr(returnCode int, err error) CliErr {
	return CliErr{ReturnCode: returnCode, Err: err}
}

type CliErr struct {
	ReturnCode int
	Err        error
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
	if !c.app.IsRepoInitialised() {
		return c.error(1, ErrUninitialisedRepo)
	}

	return nil
}

func (c *CobraHarness) ConfigList(cmd *cobra.Command, args []string) error {
	//result, err := c.app.Developer.ListConfiguration()
	//if err != nil {
	//	return NewCliErr(1, err)
	//}

	//return c.output(result)
	return nil
}
