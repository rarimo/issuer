// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	data "gitlab.com/q-dev/q-id/issuer/internal/data"
)

// ClaimsQ is an autogenerated mock type for the ClaimsQ type
type ClaimsQ struct {
	mock.Mock
}

// Get provides a mock function with given fields: id
func (_m *ClaimsQ) Get(id uint64) (*data.Claim, error) {
	ret := _m.Called(id)

	var r0 *data.Claim
	if rf, ok := ret.Get(0).(func(uint64) *data.Claim); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.Claim)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAuthClaim provides a mock function with given fields:
func (_m *ClaimsQ) GetAuthClaim() (*data.Claim, error) {
	ret := _m.Called()

	var r0 *data.Claim
	if rf, ok := ret.Get(0).(func() *data.Claim); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.Claim)
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

// GetBySchemaType provides a mock function with given fields: schemaType, userID
func (_m *ClaimsQ) GetBySchemaType(schemaType string, userID string) (*data.Claim, error) {
	ret := _m.Called(schemaType, userID)

	var r0 *data.Claim
	if rf, ok := ret.Get(0).(func(string, string) *data.Claim); ok {
		r0 = rf(schemaType, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*data.Claim)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(schemaType, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0
func (_m *ClaimsQ) Insert(_a0 *data.Claim) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*data.Claim) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// New provides a mock function with given fields:
func (_m *ClaimsQ) New() data.ClaimsQ {
	ret := _m.Called()

	var r0 data.ClaimsQ
	if rf, ok := ret.Get(0).(func() data.ClaimsQ); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(data.ClaimsQ)
		}
	}

	return r0
}

type mockConstructorTestingTNewClaimsQ interface {
	mock.TestingT
	Cleanup(func())
}

// NewClaimsQ creates a new instance of ClaimsQ. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClaimsQ(t mockConstructorTestingTNewClaimsQ) *ClaimsQ {
	mock := &ClaimsQ{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
