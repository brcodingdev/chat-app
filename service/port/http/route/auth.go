package route

import (
	"github.com/brcodingdev/chat-app/service/port/api"
	"github.com/brcodingdev/chat-app/service/port/http/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

// AuthProvider ...
type AuthProvider struct {
	router *mux.Router
	api    *api.AuthAPI
}

// NewAuthProvider ...
func NewAuthProvider(
	router *mux.Router,
	api *api.AuthAPI,
) *AuthProvider {
	return &AuthProvider{
		router: router,
		api:    api,
	}
}

// RegisterRoutes ...
func (a *AuthProvider) RegisterRoutes() {
	subRouter := a.router.PathPrefix("/api/auth").Subrouter()
	subRouter.Use(middleware.HeaderMiddleware)
	subRouter.HandleFunc("/login", a.api.Login).Methods(http.MethodPost)
	subRouter.HandleFunc("/signup", a.api.SignUp).Methods(http.MethodPost)
}
