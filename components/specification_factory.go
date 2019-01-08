package components

import (
	"fmt"
	"io"

	"github.com/endiangroup/specstack/specification"
	"github.com/spf13/afero"
)

type SpecificationFactory struct {
	FileSystem  afero.Fs
	FeaturesDir string
	WarningPipe io.Writer
}

func NewSpecificationFactory(
	fileSystem afero.Fs,
	featuresDir string,
	warningPipe io.Writer,
) *SpecificationFactory {
	return &SpecificationFactory{
		FileSystem:  fileSystem,
		FeaturesDir: featuresDir,
		WarningPipe: warningPipe,
	}
}

func (s *SpecificationFactory) SpecificationReader() specification.Reader {
	return specification.NewFilesystemReader(s.FileSystem, s.FeaturesDir)
}

func (s *SpecificationFactory) EmitWarning(warning error) {
	fmt.Fprintf(s.WarningPipe, "WARNING: %s\n", warning.Error())
}

func (s *SpecificationFactory) Specification() (*specification.Specification, specification.Reader, error) {
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
