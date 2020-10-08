package cmd

import (
	"bytes"
	"errors"
	"testing"

	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/repository"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type stdInOutErr struct {
	stdin, stdout, stderr bytes.Buffer
}

func Test_PersistenPreRunE_ReturnsRepoError(t *testing.T) {
	mockConfigAsserter := &specstack.MockConfigAsserter{}
	mockRepo := &repository.MockRepository{}
	app := &specstack.Application{
		ConfigAsserter: mockConfigAsserter,
		Repository:     mockRepo,
	}
	h, _ := setupHarness(app)

	mockRepo.On("IsInitialised").Return(false)

	err := h.PersistentPreRunE(&cobra.Command{}, nil)

	assert.Equal(t, err, NewCliErr(1, errors.New("Please initialise repository first before running")))
}

func Test_PersistenPreRunE_ReturnsConfigAssertionError(t *testing.T) {
	mockConfigAsserter := &specstack.MockConfigAsserter{}
	mockRepo := &repository.MockRepository{}
	app := &specstack.Application{
		ConfigAsserter: mockConfigAsserter,
		Repository:     mockRepo,
	}
	h, _ := setupHarness(app)

	mockRepo.On("IsInitialised").Return(true)
	mockConfigAsserter.On("AssertConfig").Return(errors.New("!!!"))

	err := h.PersistentPreRunE(&cobra.Command{}, nil)

	assert.Equal(t, err, NewCliErr(1, errors.New("!!!")))
}

func setupHarness(mockSs *specstack.Application) (*CobraHarness, *stdInOutErr) {
	io := stdInOutErr{}
	return NewCobraHarness(mockSs, &io.stdin, &io.stdout, &io.stderr), &io
}
