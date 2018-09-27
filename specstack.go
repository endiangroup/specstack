package specstack

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
)

func NewCliApp(fs afero.Fs) CliApp {
	return CliApp{
		Fs: fs,
	}
}

type CliApp struct {
	Fs afero.Fs
}

func (c CliApp) Run(args []string, stdout, stdin, stderr io.Writer) int {
	fmt.Fprintf(stderr, "Error: Please initialise git first before running: %s\n", args[0])
	return 1
}
