package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/cucumber/gherkin-go"
	"github.com/cucumber/pretty-formatter-go"
)

var dialectFlag = flag.String("dialect", "en", "Gherkin Dialect")
var writeFlag = flag.Bool("w", false, "write result to (source) file instead of stdout")

func main() {
	flag.Parse()
	paths := flag.Args()
	output, shutdown := outputStream(paths[0])
	defer func() {
		output.Flush()
		shutdown()
	}()

	if len(paths) == 0 {
		// Results mode. Read messages from STDIN
		pretty.ProcessMessages(os.Stdin, output, true)
	} else {
		// Pretty formatting mode.
		buf := &bytes.Buffer{}
		_, err := gherkin.Messages(
			[]string{paths[0]},
			nil,
			*dialectFlag,
			true,
			true,
			true,
			buf,
			false,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse Gherkin: %+v\n", err)
			os.Exit(1)
		}
		pretty.ProcessMessages(buf, output, false)
	}
}

func outputStream(finalPath string) (stream *bufio.Writer, shutdown func()) {
	if !*writeFlag {
		return bufio.NewWriter(os.Stdout), func() {}
	}

	tmpfile, err := ioutil.TempFile("", "specfmt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get tempfile path file: %s\n", err)
		os.Exit(1)
	}

	file, err := os.Create(tmpfile.Name())

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open output file: %s\n", err)
		os.Exit(1)
	}

	if err := file.Sync(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to synx output file: %s\n", err)
		os.Exit(1)
	}

	return bufio.NewWriter(file), func() {

		defer os.Remove(tmpfile.Name())

		if err := file.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to close output file: %s\n", err)
			os.Exit(1)
		}

		if err := copy(tmpfile.Name(), finalPath); err != nil {
			fmt.Fprintf(os.Stderr, "failed to copy temp file: %s\n", err)
			os.Exit(1)
		}
	}
}

func copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
