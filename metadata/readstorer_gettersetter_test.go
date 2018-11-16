package metadata

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/endiangroup/snaptest"
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

func Test_AValidGetterSetterReadStorerCanAssertMetadataHeaders(t *testing.T) {
	rs := &readStorer{}

	t.Run("Asset headers on empty entry", func(t *testing.T) {
		entry := Entry{}
		require.Nil(t, rs.assertHeaders(&entry))
		require.NotEqual(t, uuid.UUID{}, entry.Id)
		require.NotEqual(t, time.Time{}, entry.Created)
	})

	t.Run("Don't headers on non-empty entry", func(t *testing.T) {
		entry := Entry{}

		uid, err := uuid.NewV4()
		require.Nil(t, err)
		entry.Id = uid

		now := time.Now()
		entry.Created = now

		require.Nil(t, rs.assertHeaders(&entry))
		require.Equal(t, uid, entry.Id)
		require.Equal(t, now, entry.Created)
	})
}

func Test_AValidGetterSetterReadStorerCanStore(t *testing.T) {
	key, entry := bytes.NewBuffer([]byte{}), &Entry{}

	mockGetterSetter := &MockGetterSetter{}
	mockGetterSetter.On("SetMetadata", key, entry).Return(nil)

	rs := New(mockGetterSetter)
	require.Nil(t, rs.Store(key, entry))

	require.NotEqual(t, uuid.UUID{}, entry.Id)
	require.NotEqual(t, time.Time{}, entry.Created)
}

func Test_AValidGetterSetterReadStorerCanDelete(t *testing.T) {

	key := bytes.NewBuffer([]byte{})
	uids := mockUids(t, 2)
	now := time.Now()

	entries := []*Entry{
		{
			Id:      uids[1],
			Created: now,
			Name:    "A",
			Value:   "B",
		},
	}

	mockGetterSetter := &MockGetterSetter{}
	mockGetterSetter.On("GetMetadata", key, mock.Anything).
		Run(func(args mock.Arguments) {
			input := args.Get(1).(*[]*Entry)
			*input = entries
		}).
		Return(nil)

	rs := New(mockGetterSetter)

	t.Run("Fail when there's no entry", func(t *testing.T) {
		require.Equal(t, fmt.Errorf("No entry for id %s", uids[0]), rs.Delete(key, uids[0]))
	})

	t.Run("Mark entry as deleted", func(t *testing.T) {
		deleted := &Entry{
			Id:      uids[1],
			Created: now,
			Name:    "A",
			Value:   "B",
			Status:  StatusDeleted,
		}

		mockGetterSetter.On("SetMetadata", key, deleted).Return(nil)

		require.Nil(t, rs.Delete(key, uids[1]))
	})
}

func Test_AValidGetterSetterReadStorerCanRead(t *testing.T) {

	key := bytes.NewBuffer([]byte{})
	uids := mockUids(t, 4)

	entries := []*Entry{
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

	mockGetterSetter := &MockGetterSetter{}
	mockGetterSetter.On("GetMetadata", key, mock.Anything).
		Run(func(args mock.Arguments) {
			input := args.Get(1).(*[]*Entry)
			*input = entries
		}).
		Return(nil)

	rs := New(mockGetterSetter)

	t.Run("Get merged lists", func(t *testing.T) {
		entries, err := rs.Read(key)
		require.Nil(t, err)
		snaptest.Snapshot(t, entries)
	})
}
