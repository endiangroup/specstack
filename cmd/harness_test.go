package cmd

import (
	"bytes"
	"testing"

	"github.com/endiangroup/specstack"
	"github.com/stretchr/testify/assert"
)

type stdInOutErr struct {
	stdin, stdout, stderr bytes.Buffer
}

func Test_PersistenPreRunE_PrintsErrorToStderr(t *testing.T) {
	mockSs := &specstack.MockSpecStack{}
	h, io := setupHarness(mockSs)

	mockSs.On("IsRepoInitialised").Return(false)

	h.PersistentPreRunE(nil, nil)

	assert.Contains(t, io.stderr.String(), ErrUninitialisedRepo.Error())
}

func Test_PersistenPreRunE_ReturnsErrorIfUninitialisedRepo(t *testing.T) {
	mockSs := &specstack.MockSpecStack{}
	h, _ := setupHarness(mockSs)

	mockSs.On("IsRepoInitialised").Return(false)

	err := h.PersistentPreRunE(nil, nil)

	assert.Equal(t, err, NewCliErr(1, ErrUninitialisedRepo))
}

func setupHarness(mockSs *specstack.MockSpecStack) (*CobraHarness, *stdInOutErr) {
	io := stdInOutErr{}
	return NewCobraHarness(mockSs, &io.stdin, &io.stdout, &io.stderr), &io
}
