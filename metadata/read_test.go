package metadata

import (
	"bytes"
	"testing"

	"github.com/endiangroup/snaptest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_CanRead(t *testing.T) {

	key := bytes.NewBuffer([]byte{})

	entries := []*Entry{
		{
			Name:  "A",
			Value: "1",
		},
		{
			Name:  "A",
			Value: "2",
		},
		{
			Name:  "B",
			Value: "0",
		},
		{
			Name:  "A",
			Value: "3",
		},
	}

	mockStore := &MockStorer{}
	mockStore.On("ReadAllMetadata", key, mock.Anything).
		Run(func(args mock.Arguments) {
			input := args.Get(1).(*[]*Entry)
			*input = entries
		}).
		Return(nil)

	t.Run("Get merged lists", func(t *testing.T) {
		entries, err := ReadAll(mockStore, key)
		require.Nil(t, err)
		snaptest.Snapshot(t, entries)
	})
}
