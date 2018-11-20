package metadata

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_CanDelete(t *testing.T) {

	key := bytes.NewBuffer([]byte{})
	uids := mockUids(t, 2)
	now := time.Now()

	entries := []*Entry{
		{
			Id:        uids[1],
			CreatedAt: now,
			Name:      "A",
			Value:     "B",
		},
	}

	mockStore := &MockStorer{}
	mockStore.On("ReadAllMetadata", key, mock.Anything).
		Run(func(args mock.Arguments) {
			input := args.Get(1).(*[]*Entry)
			*input = entries
		}).
		Return(nil)

	t.Run("Fail when there's no entry", func(t *testing.T) {
		require.Equal(t, fmt.Errorf("No entry for id %s", uids[0]), Delete(mockStore, key, uids[0]))
	})

	t.Run("Mark entry as deleted", func(t *testing.T) {
		deleted := &Entry{
			Id:        uids[1],
			CreatedAt: now,
			Name:      "A",
			Value:     "B",
		}

		mockStore.On("StoreMetadata", key, mock.Anything).
			Run(func(args mock.Arguments) {
				output := (args.Get(1).(*Entry))
				require.NotEqual(t, time.Time{}, output.DeletedAt)

				output.DeletedAt = time.Time{}
				require.Equal(t, deleted, output)
			}).
			Return(nil)

		require.Nil(t, Delete(mockStore, key, uids[1]))
	})
}
