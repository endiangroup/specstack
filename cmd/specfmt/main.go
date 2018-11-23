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
	pretty "github.com/endiangroup/pretty-formatter-go"
)

var dialectFlag = flag.String("dialect", "en", "Gherkin Dialect")
var writeFlag = flag.Bool("w", false, "write result to (source) file instead of stdout")
var lintFlag = flag.Bool("l", false, "list files whoe formatting differes from specfmt's")

func main() {
	flag.Parse()
	paths := flag.Args()
	path := ""

	if len(paths) > 0 {
		path = paths[0]
	}
	output, shutdown := outputStream(path)
	defer func() {
		output.Flush()
		shutdown()
	}()

	switch {
	case len(paths) == 0:
		// Results mode. Read messages from STDIN
		pretty.ProcessMessages(os.Stdin, output, true)

	case *lintFlag:
		// Linting mode. List files that don't match.
		for _, path := range paths {
			invalid := false
			if !assertFileLint(path) {
				fmt.Fprintln(os.Stderr, path)
				invalid = true
			}

			if invalid {
				os.Exit(1)
			}
		}

	default:
		// Pretty formatting mode.
		buf := loadFeatureFile(paths...)
		pretty.ProcessMessages(buf, output, false)
	}
}

func loadFeatureFile(names ...string) *bytes.Buffer {
	buf := &bytes.Buffer{}
	if _, err := gherkin.Messages(
		names,
		nil,
		*dialectFlag,
		true,
		true,
		true,
		buf,
		false,
	); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse Gherkin: %+v\n", err)
		os.Exit(1)
	}
	return buf
}

func assertFileLint(path string) bool {
	input := loadFeatureFile(path)
	buf := &bytes.Buffer{}
	feature, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load Gherkin: %+v\n", err)
		os.Exit(1)
	}
	pretty.ProcessMessages(input, buf, false)
	return bytes.Equal(buf.Bytes(), feature)
}

func outputStream(finalPath string) (stream *bufio.Writer, shutdown func()) {
	if !*writeFlag || finalPath == "" {
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
