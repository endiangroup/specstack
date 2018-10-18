// Code generated by mockery v1.0.0. DO NOT EDIT.

package personas

import context "context"
import mock "github.com/stretchr/testify/mock"

// MockDeveloper is an autogenerated mock type for the Developer type
type MockDeveloper struct {
	mock.Mock
}

// GetConfiguration provides a mock function with given fields: _a0, _a1
func (_m *MockDeveloper) GetConfiguration(_a0 context.Context, _a1 string) (string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ListConfiguration provides a mock function with given fields: _a0
func (_m *MockDeveloper) ListConfiguration(_a0 context.Context) (map[string]string, error) {
	ret := _m.Called(_a0)

	var r0 map[string]string
	if rf, ok := ret.Get(0).(func(context.Context) map[string]string); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
