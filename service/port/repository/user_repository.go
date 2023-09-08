package repository

import (
	"github.com/brcodingdev/chat-app/service/pkg/model"
	"gorm.io/gorm"
)

// User ...
type User interface {
	// Add creates new user
	Add(user *model.User) error
	// GetUserByEmail get user by email
	GetUserByEmail(email string) (model.User, error)
}

// UserDB ...
type UserDB struct {
	db *gorm.DB
}

// NewUserDB ...
func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{
		db: db,
	}
}

// Add ...
func (u *UserDB) Add(user *model.User) error {
	db := u.db.Create(user)
	return ErrorCheck(db)
}

// GetUserByEmail ...
func (u *UserDB) GetUserByEmail(email string) (model.User, error) {
	var user model.User
	db := u.db.Where("Email=?", email).Find(&user)
	return user, ErrorCheck(db)
}
