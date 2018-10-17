// Code generated by mockery v1.0.0. DO NOT EDIT.

package config

import mock "github.com/stretchr/testify/mock"

// MockStorer is an autogenerated mock type for the Storer type
type MockStorer struct {
	mock.Mock
}

// CreateConfig provides a mock function with given fields: _a0
func (_m *MockStorer) CreateConfig(_a0 *Config) (*Config, error) {
	ret := _m.Called(_a0)

	var r0 *Config
	if rf, ok := ret.Get(0).(func(*Config) *Config); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Config)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*Config) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadConfig provides a mock function with given fields:
func (_m *MockStorer) LoadConfig() (*Config, error) {
	ret := _m.Called()

	var r0 *Config
	if rf, ok := ret.Get(0).(func() *Config); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*Config)
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
