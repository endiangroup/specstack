package specification

import (
	"fmt"
	"io"

	"github.com/spf13/afero"
)

type Factory struct {
	FileSystem  afero.Fs
	FeaturesDir string
	WarningPipe io.Writer
}

func NewFactory(
	fileSystem afero.Fs,
	featuresDir string,
	warningPipe io.Writer,
) *Factory {
	return &Factory{
		FileSystem:  fileSystem,
		FeaturesDir: featuresDir,
		WarningPipe: warningPipe,
	}
}

func (s *Factory) SpecificationReader() Reader {
	return NewFilesystemReader(s.FileSystem, s.FeaturesDir)
}

func (s *Factory) EmitWarning(warning error) {
	fmt.Fprintf(s.WarningPipe, "WARNING: %s\n", warning.Error())
}

func (s *Factory) Specification() (*Specification, Reader, error) {
	reader := s.SpecificationReader()

	spec, warnings, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	for _, warning := range warnings {
		s.EmitWarning(warning)
	}
	return spec, reader, nil
}
