// Code generated by mockery v2.32.3. DO NOT EDIT.

package mocks

import (
	"github.com/brcodingdev/chat-app/service/pkg/model"
	mock "github.com/stretchr/testify/mock"
)

// Chat is an autogenerated mock type for the Chat type
type Chat struct {
	mock.Mock
}

// Add provides a mock function with given fields: chat
func (_m *Chat) Add(chat *model.Chat) error {
	ret := _m.Called(chat)

	var r0 error
	if rf, ok := ret.Get(0).(func(*model.Chat) error); ok {
		r0 = rf(chat)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// List provides a mock function with given fields: roomID
func (_m *Chat) List(roomID uint) ([]model.Chat, error) {
	ret := _m.Called(roomID)

	var r0 []model.Chat
	var r1 error
	if rf, ok := ret.Get(0).(func(uint) ([]model.Chat, error)); ok {
		return rf(roomID)
	}
	if rf, ok := ret.Get(0).(func(uint) []model.Chat); ok {
		r0 = rf(roomID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Chat)
		}
	}

	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(roomID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewChat creates a new instance of Chat. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChat(t interface {
	mock.TestingT
	Cleanup(func())
}) *Chat {
	mock := &Chat{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
