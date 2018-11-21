package metadata

import (
	"bytes"
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/require"
)

func Test_APlaintextPrintScannerCanPrint(t *testing.T) {
	entries := []*Entry{
		{
			Name:  "Name A",
			Value: "B",
		},
		{
			Name:  "Longer name B",
			Value: "B",
		},
		{
			Name:  "Name C",
			Value: "Vivamus id bibendum risus: Maecenas quis arcu non ipsum bibendum posuere. Sed vitae egestas erat. Mauris amet.", //nolint:lll
		},
	}

	printer := NewPlaintextPrintscanner()
	buf := &bytes.Buffer{}

	t.Run("Can write", func(t *testing.T) {
		require.Nil(t, printer.Print(buf, entries))
		snaptest.Snapshot(t, buf.String())
	})

	t.Run("Can read", func(t *testing.T) {
		entries, err := printer.Scan(buf)
		require.Nil(t, err)
		snaptest.Snapshot(t, entries)
	})
}
