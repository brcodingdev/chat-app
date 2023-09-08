package app_test

import (
	"errors"
	"github.com/brcodingdev/chat-app/service/pkg/app"
	"github.com/brcodingdev/chat-app/service/pkg/model"
	"github.com/brcodingdev/chat-app/service/port/http/response"
	"github.com/brcodingdev/chat-app/service/port/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestListChatRoomMessages(t *testing.T) {
	cases := map[string]struct {
		roomId            uint
		modelChatMessages []model.Chat
	}{
		"list_chat_room_messages_success": {
			roomId: 1,
			modelChatMessages: []model.Chat{
				{
					UserId:     1,
					Message:    "test1",
					ChatRoomId: 1,
					User: model.User{
						Email: "clebersonh@yahoo.com.br",
					},
					ChatRoom: model.ChatRoom{
						Name: "room 1",
					},
				},
			},
		},
	}

	for caseTitle, tc := range cases {
		t.Run(caseTitle, func(t *testing.T) {
			chatRoomRepository := mocks.ChatRoom{}
			chatRepository := mocks.Chat{}
			chatRepository.On(
				"List", tc.roomId).
				Return(tc.modelChatMessages, nil).Times(1)

			chatApp := app.NewChatApp(
				&chatRoomRepository,
				&chatRepository,
			)

			res, err := chatApp.ListChatRoomMessages(tc.roomId)
			assert.NoError(t, err)
			assert.Equal(t, len(tc.modelChatMessages), len(res.Chats))
		})
	}
}

func TestListChatRoomMessagesDBError(t *testing.T) {
	cases := map[string]struct {
		roomId            uint
		modelChatMessages []model.Chat
		errorExpected     error
	}{
		"list_chat_room_messages_db_error": {
			roomId:            1,
			modelChatMessages: []model.Chat{},
			errorExpected:     errors.New("could not find chat room messages"),
		},
	}

	for caseTitle, tc := range cases {
		t.Run(caseTitle, func(t *testing.T) {
			chatRoomRepository := mocks.ChatRoom{}
			chatRepository := mocks.Chat{}
			chatRepository.On(
				"List", tc.roomId).
				Return(tc.modelChatMessages, tc.errorExpected).Times(1)

			chatApp := app.NewChatApp(
				&chatRoomRepository,
				&chatRepository,
			)

			res, err := chatApp.ListChatRoomMessages(tc.roomId)
			assert.Error(t, err)
			assert.Equal(t, response.ChatRoomMessagesResponse{}, res)
		})
	}
}
