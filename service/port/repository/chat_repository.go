package repository

import (
	"github.com/brcodingdev/chat-app/service/pkg/model"
	"gorm.io/gorm"
)

// Chat ...
type Chat interface {
	// Add create new chat
	Add(chat *model.Chat) error
	// List lists chats limited by 50 rows
	List(roomID uint) ([]model.Chat, error)
}

// ChatDB ...
type ChatDB struct {
	db *gorm.DB
}

// NewChatDB ...
func NewChatDB(db *gorm.DB) *ChatDB {
	return &ChatDB{db: db}
}

// Add ...
func (c *ChatDB) Add(chat *model.Chat) error {
	db := c.db.Where(chat).FirstOrCreate(chat)
	return ErrorCheck(db)
}

// List ...
func (c *ChatDB) List(roomID uint) ([]model.Chat, error) {
	var chatList []model.Chat
	db := c.db.Where(
		model.Chat{
			ChatRoomId: roomID,
		}).Preload("ChatRoom").Preload("User").Find(
		&chatList,
	).Limit(50)

	return chatList, ErrorCheck(db)
}
