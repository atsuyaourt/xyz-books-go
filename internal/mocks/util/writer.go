// Code generated by mockery v2.34.2. DO NOT EDIT.

package mockutil

import mock "github.com/stretchr/testify/mock"

// MockWriter is an autogenerated mock type for the Writer type
type MockWriter struct {
	mock.Mock
}

type MockWriter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockWriter) EXPECT() *MockWriter_Expecter {
	return &MockWriter_Expecter{mock: &_m.Mock}
}

// Error provides a mock function with given fields:
func (_m *MockWriter) Error() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockWriter_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type MockWriter_Error_Call struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
func (_e *MockWriter_Expecter) Error() *MockWriter_Error_Call {
	return &MockWriter_Error_Call{Call: _e.mock.On("Error")}
}

func (_c *MockWriter_Error_Call) Run(run func()) *MockWriter_Error_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockWriter_Error_Call) Return(_a0 error) *MockWriter_Error_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWriter_Error_Call) RunAndReturn(run func() error) *MockWriter_Error_Call {
	_c.Call.Return(run)
	return _c
}

// Flush provides a mock function with given fields:
func (_m *MockWriter) Flush() {
	_m.Called()
}

// MockWriter_Flush_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Flush'
type MockWriter_Flush_Call struct {
	*mock.Call
}

// Flush is a helper method to define mock.On call
func (_e *MockWriter_Expecter) Flush() *MockWriter_Flush_Call {
	return &MockWriter_Flush_Call{Call: _e.mock.On("Flush")}
}

func (_c *MockWriter_Flush_Call) Run(run func()) *MockWriter_Flush_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockWriter_Flush_Call) Return() *MockWriter_Flush_Call {
	_c.Call.Return()
	return _c
}

func (_c *MockWriter_Flush_Call) RunAndReturn(run func()) *MockWriter_Flush_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: record
func (_m *MockWriter) Write(record []string) error {
	ret := _m.Called(record)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(record)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockWriter_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type MockWriter_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - record []string
func (_e *MockWriter_Expecter) Write(record interface{}) *MockWriter_Write_Call {
	return &MockWriter_Write_Call{Call: _e.mock.On("Write", record)}
}

func (_c *MockWriter_Write_Call) Run(run func(record []string)) *MockWriter_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]string))
	})
	return _c
}

func (_c *MockWriter_Write_Call) Return(_a0 error) *MockWriter_Write_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockWriter_Write_Call) RunAndReturn(run func([]string) error) *MockWriter_Write_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockWriter creates a new instance of MockWriter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockWriter(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockWriter {
	mock := &MockWriter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
