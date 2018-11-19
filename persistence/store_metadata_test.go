package persistence

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/endiangroup/snaptest"
	"github.com/endiangroup/specstack/metadata"
	uuid "github.com/satori/go.uuid"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mockUids(t *testing.T, count int) (output []uuid.UUID) {
	for i := 0; i < count; i++ {
		uid := uuid.NewV5(uuid.NamespaceOID, fmt.Sprintf("%d", i))
		output = append(output, uid)
	}
	return
}

func Test_StoreMetadata_CanAssertMetadataHeaders(t *testing.T) {
	mockConfigStorer := &MockConfigStorer{}
	mockMetadataStorer := &MockMetadataStorer{}
	rs := NewStore(mockConfigStorer, mockMetadataStorer)

	t.Run("Assert headers on empty entry", func(t *testing.T) {
		entry := metadata.Entry{}
		require.Nil(t, rs.assertHeaders(&entry))
		require.NotEqual(t, uuid.UUID{}, entry.Id)
		require.NotEqual(t, time.Time{}, entry.Created)
	})

	t.Run("Don't populate headers on non-empty entry", func(t *testing.T) {
		entry := metadata.Entry{}

		uid := uuid.NewV4()
		entry.Id = uid

		now := time.Now()
		entry.Created = now

		require.Nil(t, rs.assertHeaders(&entry))
		require.Equal(t, uid, entry.Id)
		require.Equal(t, now, entry.Created)
	})
}

func Test_StoreMetadata_CanStoreValueData(t *testing.T) {
	key, entry := bytes.NewBuffer([]byte{}), &metadata.Entry{}

	mockMetadataStore := &MockMetadataStorer{}
	mockMetadataStore.On("SetMetadata", key, entry).Return(nil)

	mockConfigStorer := &MockConfigStorer{}
	rs := NewStore(mockConfigStorer, mockMetadataStore)

	require.Nil(t, rs.StoreMetadata(key, entry))

	require.NotEqual(t, uuid.UUID{}, entry.Id)
	require.NotEqual(t, time.Time{}, entry.Created)
}

func Test_StoreMetadata_CanDelete(t *testing.T) {

	key := bytes.NewBuffer([]byte{})
	uids := mockUids(t, 2)
	now := time.Now()

	entries := []*metadata.Entry{
		{
			Id:      uids[1],
			Created: now,
			Name:    "A",
			Value:   "B",
		},
	}

	mockMetadataStore := &MockMetadataStorer{}
	mockMetadataStore.On("GetMetadata", key, mock.Anything).
		Run(func(args mock.Arguments) {
			input := args.Get(1).(*[]*metadata.Entry)
			*input = entries
		}).
		Return(nil)

	mockConfigStorer := &MockConfigStorer{}
	rs := NewStore(mockConfigStorer, mockMetadataStore)

	t.Run("Fail when there's no entry", func(t *testing.T) {
		require.Equal(t, fmt.Errorf("No entry for id %s", uids[0]), rs.DeleteMetadata(key, uids[0]))
	})

	t.Run("Mark entry as deleted", func(t *testing.T) {
		deleted := &metadata.Entry{
			Id:      uids[1],
			Created: now,
			Name:    "A",
			Value:   "B",
			Status:  metadata.StatusDeleted,
		}

		mockMetadataStore.On("SetMetadata", key, deleted).Return(nil)

		require.Nil(t, rs.DeleteMetadata(key, uids[1]))
	})
}

func Test_StoreMetadata_CanRead(t *testing.T) {

	key := bytes.NewBuffer([]byte{})
	uids := mockUids(t, 4)

	entries := []*metadata.Entry{
		{
			Id:    uids[0],
			Name:  "A",
			Value: "1",
		},
		{
			Id:    uids[1],
			Name:  "A",
			Value: "2",
		},
		{
			Id:    uids[3],
			Name:  "B",
			Value: "0",
		},
		{
			Id:    uids[2],
			Name:  "A",
			Value: "3",
		},
	}

	mockMetadataStore := &MockMetadataStorer{}
	mockMetadataStore.On("GetMetadata", key, mock.Anything).
		Run(func(args mock.Arguments) {
			input := args.Get(1).(*[]*metadata.Entry)
			*input = entries
		}).
		Return(nil)

	mockConfigStorer := &MockConfigStorer{}
	rs := NewStore(mockConfigStorer, mockMetadataStore)

	t.Run("Get merged lists", func(t *testing.T) {
		entries, err := rs.ReadMetadata(key)
		require.Nil(t, err)
		snaptest.Snapshot(t, entries)
	})
}
