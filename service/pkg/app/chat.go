package app

import (
	"github.com/brcodingdev/chat-app/service/pkg/model"
	"github.com/brcodingdev/chat-app/service/port/http/response"
	"github.com/brcodingdev/chat-app/service/port/repository"
)

// Chat ...
type Chat interface {
	// CreateChatRoom creates a room
	CreateChatRoom(name string) (response.ChatRoomResponse, error)
	// CreateChatMessage creates a message
	CreateChatMessage(msg string, roomId, userId uint) bool
	// ListChatRooms lists rooms
	ListChatRooms() (response.ChatRoomsResponse, error)
	// ListChatRoomMessages lists last messages
	ListChatRoomMessages(roomId uint) (response.ChatRoomMessagesResponse, error)
}

// ChatApp ...
type ChatApp struct {
	chatRoomRepository repository.ChatRoom
	chatRepository     repository.Chat
}

// NewChatApp ...
func NewChatApp(
	chatRoomRepository repository.ChatRoom,
	chatRepository repository.Chat,
) *ChatApp {
	return &ChatApp{
		chatRoomRepository: chatRoomRepository,
		chatRepository:     chatRepository,
	}
}

// ListChatRooms ...
func (c *ChatApp) ListChatRooms() (response.ChatRoomsResponse, error) {
	var chtList []model.ChatRoom
	chtList, err := c.chatRoomRepository.List()
	if err != nil {
		return response.ChatRoomsResponse{},
			err
	}

	return response.ChatRoomsResponse{ChatRooms: chtList}, nil
}

// CreateChatRoom ...
func (c *ChatApp) CreateChatRoom(
	name string,
) (response.ChatRoomResponse, error) {
	chatRoom := model.ChatRoom{
		Name: name,
	}
	err := c.chatRoomRepository.Add(&chatRoom)

	if err != nil {
		return response.ChatRoomResponse{}, err
	}
	return response.ChatRoomResponse{ChatRoom: chatRoom}, nil
}

// ListChatRoomMessages ...
func (c *ChatApp) ListChatRoomMessages(
	roomId uint,
) (response.ChatRoomMessagesResponse, error) {
	var chtMsgList []response.ChatMessage
	chtList, err := c.chatRepository.List(roomId)

	if err != nil {
		return response.ChatRoomMessagesResponse{}, err
	}

	for _, v := range chtList {
		chtMsgList = append(chtMsgList, response.ChatMessage{
			ChatMessage:  v.Message,
			ChatUser:     v.User.Email,
			ChatRoomId:   v.ChatRoomId,
			ChatRoomName: v.ChatRoom.Name,
		})
	}
	return response.ChatRoomMessagesResponse{Chats: chtMsgList}, nil
}

// CreateChatMessage ...
func (c *ChatApp) CreateChatMessage(msg string, roomId, userId uint) bool {
	chat := model.Chat{
		Message:    msg,
		UserId:     userId,
		ChatRoomId: roomId,
	}

	err := c.chatRepository.Add(&chat)
	if err != nil {

		return false
	}

	return true
}
