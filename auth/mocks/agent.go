// Code generated by mockery v2.43.2. DO NOT EDIT.

// Copyright (c) Abstract Machines

package mocks

import (
	context "context"

	auth "github.com/absmach/magistrala/auth"

	mock "github.com/stretchr/testify/mock"
)

// PolicyAgent is an autogenerated mock type for the PolicyAgent type
type PolicyAgent struct {
	mock.Mock
}

// CheckPolicy provides a mock function with given fields: ctx, pr
func (_m *PolicyAgent) CheckPolicy(ctx context.Context, pr auth.PolicyReq) error {
	ret := _m.Called(ctx, pr)

	if len(ret) == 0 {
		panic("no return value specified for CheckPolicy")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, auth.PolicyReq) error); ok {
		r0 = rf(ctx, pr)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewPolicyAgent creates a new instance of PolicyAgent. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPolicyAgent(t interface {
	mock.TestingT
	Cleanup(func())
}) *PolicyAgent {
	mock := &PolicyAgent{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
