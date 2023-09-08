package response

import (
	"github.com/brcodingdev/chat-app/service/pkg/model"
)

// ChatRoomsResponse ...
type ChatRoomsResponse struct {
	ChatRooms []model.ChatRoom `json:"ChatRooms"`
}

// ChatRoomResponse ...
type ChatRoomResponse struct {
	ChatRoom model.ChatRoom `json:"ChatRoom"`
}

// ChatMessage ...
type ChatMessage struct {
	ChatMessage  string `json:"chatMessage"`
	ChatUser     string `json:"chatUser"`
	ChatRoomId   uint   `json:"chatRoomId"`
	ChatRoomName string `json:"chatRoomName"`
}

// ChatRoomMessagesResponse ...
type ChatRoomMessagesResponse struct {
	Chats []ChatMessage `json:"Chats"`
}
