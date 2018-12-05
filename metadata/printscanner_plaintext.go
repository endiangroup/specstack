package metadata

import (
	"bufio"
	"fmt"
	io "io"
	"strings"
)

const lineLength = 100

type PlaintextPrintScanner struct {
}

func NewPlaintextPrintscanner() PrintScanner {
	return &PlaintextPrintScanner{}
}

func (p *PlaintextPrintScanner) Print(writer io.Writer, entries []*Entry) error {
	longest := 0

	for _, entry := range entries {
		if l := len(entry.Name); l > longest {
			longest = l
		}
	}

	for _, entry := range entries {
		if _, err := fmt.Fprintf(
			writer,
			"%s: %s\n",
			p.padRight(entry.Name, " ", longest),
			p.truncate(entry.Value, lineLength-longest-2),
		); err != nil {
			return err
		}
	}

	return nil
}

func (p *PlaintextPrintScanner) Scan(reader io.Reader) ([]Entry, error) {
	entries := []Entry{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")

		if len(parts) > 1 {
			e := Entry{
				Name:  strings.TrimSpace(parts[0]),
				Value: strings.TrimSpace(strings.Join(parts[1:], ":")),
			}
			entries = append(entries, e)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (p *PlaintextPrintScanner) padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}

func (p *PlaintextPrintScanner) truncate(str string, num int) string {
	output := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		output = str[0:num] + "..."
	}
	return output
}
