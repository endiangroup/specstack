// Code generated by mockery v1.0.0. DO NOT EDIT.

package repository

import io "io"
import mock "github.com/stretchr/testify/mock"

// MockRepository is an autogenerated mock type for the Repository type
type MockRepository struct {
	mock.Mock
}

// AllConfig provides a mock function with given fields:
func (_m *MockRepository) AllConfig() (map[string]string, error) {
	ret := _m.Called()

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func() map[string]string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConfig provides a mock function with given fields: _a0
func (_m *MockRepository) GetConfig(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetMetadata provides a mock function with given fields: key, output
func (_m *MockRepository) GetMetadata(key io.Reader, output interface{}) error {
	ret := _m.Called(key, output)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Reader, interface{}) error); ok {
		r0 = rf(key, output)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Init provides a mock function with given fields:
func (_m *MockRepository) Init() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// IsInitialised provides a mock function with given fields:
func (_m *MockRepository) IsInitialised() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// SetConfig provides a mock function with given fields: _a0, _a1
func (_m *MockRepository) SetConfig(_a0 string, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetMetadata provides a mock function with given fields: key, value
func (_m *MockRepository) SetMetadata(key io.Reader, value interface{}) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(io.Reader, interface{}) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnsetConfig provides a mock function with given fields: _a0
func (_m *MockRepository) UnsetConfig(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
