// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	data "gitlab.com/q-dev/q-id/issuer/internal/data"
)

// TreeStorageQ is an autogenerated mock type for the TreeStorageQ type
type TreeStorageQ struct {
	mock.Mock
}

// Get provides a mock function with given fields: key
func (_m *TreeStorageQ) Get(key []byte) ([]byte, error) {
	ret := _m.Called(key)

	var r0 []byte
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(key)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(key)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: key, value
func (_m *TreeStorageQ) Insert(key []byte, value []byte) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []byte) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// New provides a mock function with given fields:
func (_m *TreeStorageQ) New() data.TreeStorageQ {
	ret := _m.Called()

	var r0 data.TreeStorageQ
	if rf, ok := ret.Get(0).(func() data.TreeStorageQ); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(data.TreeStorageQ)
		}
	}

	return r0
}

// Upsert provides a mock function with given fields: key, value
func (_m *TreeStorageQ) Upsert(key []byte, value []byte) error {
	ret := _m.Called(key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, []byte) error); ok {
		r0 = rf(key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewTreeStorageQ interface {
	mock.TestingT
	Cleanup(func())
}

// NewTreeStorageQ creates a new instance of TreeStorageQ. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTreeStorageQ(t mockConstructorTestingTNewTreeStorageQ) *TreeStorageQ {
	mock := &TreeStorageQ{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
