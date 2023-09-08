package model

import (
	"gorm.io/gorm"
)

// ChatRoom ...
type ChatRoom struct {
	gorm.Model
	Name string `json:"Name"`
}
