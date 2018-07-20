// Code generated by mockery v1.0.0. DO NOT EDIT.
package net

import (
	context "context"
	"net"

	mock "github.com/stretchr/testify/mock"
)

// MockMessenger is an autogenerated mock type for the Messenger type
type MockMessenger struct {
	mock.Mock
}

// Handle provides a mock function with given fields: contentType, h
func (_m *MockMessenger) Handle(contentType string, h EnvelopeHandler) error {
	ret := _m.Called(contentType, h)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, EnvelopeHandler) error); ok {
		r0 = rf(contentType, h)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Listen provides a mock function with given fields: ctx, addrress
func (_m *MockMessenger) Listen(ctx context.Context, addrress string) (net.Listener, error) {
	ret := _m.Called(ctx, addrress)

	var r0 net.Listener
	if rf, ok := ret.Get(0).(func(context.Context, string) net.Listener); ok {
		r0 = rf(ctx, addrress)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Listener)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, addrress)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Send provides a mock function with given fields: ctx, envelope
func (_m *MockMessenger) Send(ctx context.Context, envelope *Envelope) error {
	ret := _m.Called(ctx, envelope)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *Envelope) error); ok {
		r0 = rf(ctx, envelope)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}