// Code generated by mockery v1.0.0. DO NOT EDIT.
package net

import mock "github.com/stretchr/testify/mock"

// MockAddressBook is an autogenerated mock type for the AddressBook type
type MockAddressBook struct {
	mock.Mock
}

// CreateNewPeer provides a mock function with given fields:
func (_m *MockAddressBook) CreateNewPeer() (*PrivatePeerInfo, error) {
	ret := _m.Called()

	var r0 *PrivatePeerInfo
	if rf, ok := ret.Get(0).(func() *PrivatePeerInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PrivatePeerInfo)
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

// GetAllPeerInfo provides a mock function with given fields:
func (_m *MockAddressBook) GetAllPeerInfo() ([]*PeerInfo, error) {
	ret := _m.Called()

	var r0 []*PeerInfo
	if rf, ok := ret.Get(0).(func() []*PeerInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*PeerInfo)
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

// GetLocalPeerInfo provides a mock function with given fields:
func (_m *MockAddressBook) GetLocalPeerInfo() *PrivatePeerInfo {
	ret := _m.Called()

	var r0 *PrivatePeerInfo
	if rf, ok := ret.Get(0).(func() *PrivatePeerInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PrivatePeerInfo)
		}
	}

	return r0
}

// GetPeerInfo provides a mock function with given fields: peerID
func (_m *MockAddressBook) GetPeerInfo(peerID string) (*PeerInfo, error) {
	ret := _m.Called(peerID)

	var r0 *PeerInfo
	if rf, ok := ret.Get(0).(func(string) *PeerInfo); ok {
		r0 = rf(peerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PeerInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(peerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadOrCreateLocalPeerInfo provides a mock function with given fields: path
func (_m *MockAddressBook) LoadOrCreateLocalPeerInfo(path string) (*PrivatePeerInfo, error) {
	ret := _m.Called(path)

	var r0 *PrivatePeerInfo
	if rf, ok := ret.Get(0).(func(string) *PrivatePeerInfo); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PrivatePeerInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LoadPrivatePeerInfo provides a mock function with given fields: path
func (_m *MockAddressBook) LoadPrivatePeerInfo(path string) (*PrivatePeerInfo, error) {
	ret := _m.Called(path)

	var r0 *PrivatePeerInfo
	if rf, ok := ret.Get(0).(func(string) *PrivatePeerInfo); ok {
		r0 = rf(path)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*PrivatePeerInfo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(path)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PutLocalPeerInfo provides a mock function with given fields: _a0
func (_m *MockAddressBook) PutLocalPeerInfo(_a0 *PrivatePeerInfo) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*PrivatePeerInfo) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PutPeerInfoFromEnvelope provides a mock function with given fields: _a0
func (_m *MockAddressBook) PutPeerInfoFromEnvelope(_a0 *PeerInfo) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*PeerInfo) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorePrivatePeerInfo provides a mock function with given fields: pi, path
func (_m *MockAddressBook) StorePrivatePeerInfo(pi *PrivatePeerInfo, path string) error {
	ret := _m.Called(pi, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(*PrivatePeerInfo, string) error); ok {
		r0 = rf(pi, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}