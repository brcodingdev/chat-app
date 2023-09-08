package api

import (
	"encoding/json"
	"github.com/brcodingdev/chat-app/service/errors"
	"github.com/brcodingdev/chat-app/service/pkg/app"
	"github.com/brcodingdev/chat-app/service/port/http/request"
	"net/http"
)

// AuthAPI ...
type AuthAPI struct {
	app app.Auth
}

// NewAuthAPI ...
func NewAuthAPI(app app.Auth) *AuthAPI {
	return &AuthAPI{app: app}
}

// Login ...
func (a *AuthAPI) Login(w http.ResponseWriter, r *http.Request) {
	loginRequest := request.LoginRequest{}
	err := parseBody(r, &loginRequest)
	if err != nil {
		errResponse(errors.ErrInRequest, w)
		return
	}

	res, err := a.app.Login(loginRequest.Email, loginRequest.Password)
	if err != nil {
		errResponse(err, w)
		return
	}

	res.User.Password = ""
	data, err := json.Marshal(res)
	if err != nil {
		errResponse(err, w)
		return
	}

	ok(data, w)
}

// SignUp ...
func (a *AuthAPI) SignUp(w http.ResponseWriter, r *http.Request) {
	signUpRequest := request.SignUpRequest{}
	err := parseBody(r, &signUpRequest)
	if err != nil {
		errResponse(errors.ErrInRequest, w)
		return
	}

	res, err := a.app.SignUp(
		signUpRequest.Email,
		signUpRequest.UserName,
		signUpRequest.Password,
	)

	if err != nil {
		errResponse(err, w)
		return
	}

	res.User.Password = ""
	data, err := json.Marshal(res)
	if err != nil {
		errResponse(err, w)
		return
	}

	ok(data, w)
}
