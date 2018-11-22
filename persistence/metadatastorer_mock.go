// Code generated by mockery v1.0.0. DO NOT EDIT.

package persistence

import io "io"
import mock "github.com/stretchr/testify/mock"

// MockMetadataStorer is an autogenerated mock type for the MetadataStorer type
type MockMetadataStorer struct {
	mock.Mock
}

// GetMetadata provides a mock function with given fields: key
func (_m *MockMetadataStorer) GetMetadata(key io.Reader) ([][]byte, error) {
	ret := _m.Called(key)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func(io.Reader) [][]byte); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(io.Reader) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetMetadata provides a mock function with given fields: key, value
func (_m *MockMetadataStorer) SetMetadata(key io.Reader, value []byte) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Reader, []byte) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}