package metadata

import io "io"

type Printer interface {
	Print(io.Writer, []*Entry) error
}

type Scanner interface {
	Scan(io.Reader) ([]Entry, error)
}

type PrintScanner interface {
	Printer
	Scanner
}
