// Code generated by mockery v1.0.0. DO NOT EDIT.

package specification

import mock "github.com/stretchr/testify/mock"

// MockSourcer is an autogenerated mock type for the Sourcer type
type MockSourcer struct {
	mock.Mock
}

// Source provides a mock function with given fields:
func (_m *MockSourcer) Source() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
