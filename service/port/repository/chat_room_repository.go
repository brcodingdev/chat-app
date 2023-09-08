package repository

import (
	"github.com/brcodingdev/chat-app/service/pkg/model"
	"gorm.io/gorm"
)

// ChatRoom ...
type ChatRoom interface {
	// Add create new chat room
	Add(chatRoom *model.ChatRoom) error
	// List lists all chat rooms
	List() ([]model.ChatRoom, error)
}

// ChatRoomDB ...
type ChatRoomDB struct {
	db *gorm.DB
}

// NewChatRoomDB ...
func NewChatRoomDB(db *gorm.DB) *ChatRoomDB {
	return &ChatRoomDB{db: db}
}

// Add ...
func (c *ChatRoomDB) Add(chatRoom *model.ChatRoom) error {
	db := c.db.Where(chatRoom).FirstOrCreate(chatRoom)
	return ErrorCheck(db)
}

// List ...
func (c *ChatRoomDB) List() ([]model.ChatRoom, error) {
	var chatRooms []model.ChatRoom
	db := c.db.Order("id DESC").Find(&chatRooms)
	return chatRooms, ErrorCheck(db)
}
