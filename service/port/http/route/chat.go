package route

import (
	"github.com/brcodingdev/chat-app/service/port/api"
	"github.com/brcodingdev/chat-app/service/port/http/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

// ChatProvider ...
type ChatProvider struct {
	router *mux.Router
	api    *api.ChatAPI
}

// NewChatProvider ...
func NewChatProvider(
	router *mux.Router,
	api *api.ChatAPI,
) *ChatProvider {
	return &ChatProvider{
		router: router,
		api:    api,
	}
}

// RegisterRoutes ...
func (c *ChatProvider) RegisterRoutes() {
	subRouter := c.router.PathPrefix("/api/chat").Subrouter()
	subRouter.Use(middleware.HeaderMiddleware)
	subRouter.Use(middleware.Authenticated)

	subRouter.HandleFunc(
		"/create",
		c.api.Create,
	).Methods(http.MethodPost)

	subRouter.HandleFunc(
		"/rooms",
		c.api.ChatRooms,
	).Methods(http.MethodPost)

	subRouter.HandleFunc(
		"/room-messages",
		c.api.ChatRoomMessages,
	).Methods(http.MethodPost)
}
