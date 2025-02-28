// Code generated by mockery v2.50.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/guil95/ports-service/internal/core/domain"
	mock "github.com/stretchr/testify/mock"
)

// ParserPort is an autogenerated mock type for the ParserPort type
type ParserPort struct {
	mock.Mock
}

// Parse provides a mock function with given fields: ctx
func (_m *ParserPort) Parse(ctx context.Context) (<-chan domain.Port, <-chan error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for Parse")
	}

	var r0 <-chan domain.Port
	var r1 <-chan error
	if rf, ok := ret.Get(0).(func(context.Context) (<-chan domain.Port, <-chan error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) <-chan domain.Port); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan domain.Port)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) <-chan error); ok {
		r1 = rf(ctx)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(<-chan error)
		}
	}

	return r0, r1
}

// NewParserPort creates a new instance of ParserPort. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewParserPort(t interface {
	mock.TestingT
	Cleanup(func())
}) *ParserPort {
	mock := &ParserPort{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
