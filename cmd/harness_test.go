package cmd

import (
	"bytes"
	"testing"

	"github.com/endiangroup/specstack"
	"github.com/stretchr/testify/assert"
)

func Test_PersistenPreRunE_PrintsErrorToStderr(t *testing.T) {
	mockSs := &specstack.MockSpecStack{}
	var stdin, stdout, stderr bytes.Buffer
	h := NewCobraHarness(mockSs, &stdin, &stdout, &stderr)

	mockSs.On("IsRepoInitialised").Return(false)

	h.PersistentPreRunE(nil, nil)

	assert.Contains(t, stderr.String(), ErrUninitialisedRepo.Error())
}

func Test_PersistenPreRunE_ReturnsErrorIfUninitialisedRepo(t *testing.T) {
	mockSs := &specstack.MockSpecStack{}
	var stdin, stdout, stderr bytes.Buffer
	h := NewCobraHarness(mockSs, &stdin, &stdout, &stderr)

	mockSs.On("IsRepoInitialised").Return(false)

	err := h.PersistentPreRunE(nil, nil)

	assert.Equal(t, err, NewCliErr(1, ErrUninitialisedRepo))
}
