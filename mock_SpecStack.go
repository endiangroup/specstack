// Code generated by mockery v1.0.0. DO NOT EDIT.

package specstack

import mock "github.com/stretchr/testify/mock"

// MockSpecStack is an autogenerated mock type for the SpecStack type
type MockSpecStack struct {
	mock.Mock
}

// IsRepoInitialised provides a mock function with given fields:
func (_m *MockSpecStack) IsRepoInitialised() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}