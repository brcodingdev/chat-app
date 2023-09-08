package api

import (
	"encoding/json"
	"github.com/brcodingdev/chat-app/service/errors"
	"github.com/brcodingdev/chat-app/service/pkg/app"
	"github.com/brcodingdev/chat-app/service/port/http/request"
	"net/http"
)

// ChatAPI ...
type ChatAPI struct {
	app app.Chat
}

// NewChatAPI ...
func NewChatAPI(app app.Chat) *ChatAPI {
	return &ChatAPI{
		app: app,
	}
}

// ChatRooms ...
func (c *ChatAPI) ChatRooms(w http.ResponseWriter, r *http.Request) {
	res, err := c.app.ListChatRooms()
	if err != nil {
		errResponse(err, w)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		errResponse(err, w)
		return
	}

	ok(data, w)
}

// Create ...
func (c *ChatAPI) Create(w http.ResponseWriter, r *http.Request) {
	createRoomRequest := request.ChatRoomCreateRequest{}
	err := parseBody(r, &createRoomRequest)
	if err != nil {
		errResponse(errors.ErrInRequest, w)
		return
	}

	res, err := c.app.CreateChatRoom(createRoomRequest.Name)
	if err != nil {
		errResponse(err, w)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		errResponse(err, w)
		return
	}

	ok(data, w)
}

// ChatRoomMessages ...
func (c *ChatAPI) ChatRoomMessages(w http.ResponseWriter, r *http.Request) {
	chatMessagesRequest := request.ChatRoomMessagesRequest{}
	err := parseBody(r, &chatMessagesRequest)
	if err != nil {
		errResponse(errors.ErrInRequest, w)
		return
	}

	res, err := c.app.ListChatRoomMessages(chatMessagesRequest.RoomId)
	if err != nil {
		errResponse(err, w)
		return
	}

	data, err := json.Marshal(res)
	if err != nil {
		errResponse(err, w)
		return
	}

	ok(data, w)
}
