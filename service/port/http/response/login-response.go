package response

import (
	"github.com/brcodingdev/chat-app/service/pkg/model"
)

// LoginResponse ...
type LoginResponse struct {
	User     model.User `json:"User"`
	JwtToken string     `json:"Token"`
}
