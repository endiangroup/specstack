package persistence

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_StoreMetadata_CanStoreAnyValidValueAsMetadata(t *testing.T) {
	mockConfigStorer := &MockConfigStorer{}
	mockMetadataStorer := &MockMetadataStorer{}
	rs := NewStore(mockConfigStorer, mockMetadataStorer)

	t.Run("Throw error with invalid values", func(t *testing.T) {
		err := rs.StoreMetadata(bytes.NewBufferString(t.Name()), func() {})
		require.NotNil(t, err)
		require.Equal(t, "failed to serialise metadata value: json: unsupported type: func()", err.Error())
	})

	t.Run("Encode valid values", func(t *testing.T) {
		value := "{ / * string \n`}"
		jsn, err := json.Marshal(value)
		require.Nil(t, err)

		key := bytes.NewBufferString(t.Name())
		mockMetadataStorer.On("SetMetadata", key, jsn).Return(nil)

		require.Nil(t, rs.StoreMetadata(key, value))
	})
}

func Test_StoreMetadata_CanReadAllMetadataFromASource(t *testing.T) {
	mockConfigStorer := &MockConfigStorer{}
	mockMetadataStorer := &MockMetadataStorer{}
	rs := NewStore(mockConfigStorer, mockMetadataStorer)

	t.Run("Throw error when provided invalid values", func(t *testing.T) {
		key := bytes.NewBufferString(t.Name())
		mockMetadataStorer.On("GetMetadata", key).Return(nil, fmt.Errorf("some error"))
		require.Equal(t, fmt.Errorf("failed to get raw metadata: some error"), rs.ReadAllMetadata(key, nil))
	})

	t.Run("Retrieve valid values", func(t *testing.T) {
		key := bytes.NewBufferString(t.Name())
		type MyObject struct {
			A int
			B int
		}
		objects := []MyObject{}
		object := MyObject{1, 2}

		jsn, err := json.Marshal(object)
		require.Nil(t, err)

		mockMetadataStorer.On("GetMetadata", key).Return([][]byte{jsn}, nil)
		require.Nil(t, rs.ReadAllMetadata(key, &objects))
		require.Equal(t, []MyObject{object}, objects)
	})
}
