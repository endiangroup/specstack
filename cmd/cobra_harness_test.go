package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/endiangroup/specstack"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type stdInOutErr struct {
	stdin, stdout, stderr bytes.Buffer
}

func Test_PersistenPreRunE_ReturnsInitialiseError(t *testing.T) {
	mockSs := &specstack.MockSpecStack{}
	h, _ := setupHarness(mockSs)

	mockSs.On("Initialise").Return(errors.New("!!!"))

	err := h.PersistentPreRunE(&cobra.Command{}, nil)

	assert.Equal(t, err, NewCliErr(1, errors.New("!!!")))
}

func setupHarness(mockSs *specstack.MockSpecStack) (*CobraHarness, *stdInOutErr) {
	io := stdInOutErr{}
	return NewCobraHarness(mockSs, &io.stdin, &io.stdout, &io.stderr), &io
}
