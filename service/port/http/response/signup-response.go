package response

import (
	"github.com/brcodingdev/chat-app/service/pkg/model"
)

// SignUpResponse ...
type SignUpResponse struct {
	User     model.User `json:"User"`
	JwtToken string     `json:"Token"`
}
