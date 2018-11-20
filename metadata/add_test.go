package metadata

import (
	"bytes"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func Test_CRUDLayerCanAssertMetadataHeaders(t *testing.T) {

	t.Run("Assert headers on empty entry", func(t *testing.T) {
		entry := Entry{}
		require.Nil(t, assertHeaders(&entry))
		require.NotEqual(t, uuid.UUID{}, entry.Id)
		require.NotEqual(t, time.Time{}, entry.CreatedAt)
	})

	t.Run("Don't populate headers on non-empty entry", func(t *testing.T) {
		entry := Entry{}

		uid := uuid.NewV4()
		entry.Id = uid

		now := time.Now()
		entry.CreatedAt = now

		require.Nil(t, assertHeaders(&entry))
		require.Equal(t, uid, entry.Id)
		require.Equal(t, now, entry.CreatedAt)
	})
}

func Test_CRUDLayerCanStoreValueData(t *testing.T) {
	key, entry := bytes.NewBuffer([]byte{}), &Entry{}

	mockStorer := &MockStorer{}
	mockStorer.On("StoreMetadata", key, entry).Return(nil)

	require.Nil(t, Add(mockStorer, key, entry))

	require.NotEqual(t, uuid.UUID{}, entry.Id)
	require.NotEqual(t, time.Time{}, entry.CreatedAt)
}
