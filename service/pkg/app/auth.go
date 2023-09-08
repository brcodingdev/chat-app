package app

import (
	"github.com/brcodingdev/chat-app/service/errors"
	"github.com/brcodingdev/chat-app/service/pkg/model"
	"github.com/brcodingdev/chat-app/service/port/http/response"
	"github.com/brcodingdev/chat-app/service/port/repository"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Auth ...
type Auth interface {
	Login(userName, password string) (response.LoginResponse, error)
	SignUp(email, userName, password string) (response.SignUpResponse, error)
}

// AuthApp ...
type AuthApp struct {
	userRepository repository.User
}

// NewAuthApp ...
func NewAuthApp(userRepository repository.User) *AuthApp {
	return &AuthApp{
		userRepository: userRepository,
	}
}

// Login ...
func (a *AuthApp) Login(
	email,
	password string,
) (response.LoginResponse, error) {
	user, err := a.userRepository.GetUserByEmail(email)

	if err != nil || user.ID <= 0 {
		return response.LoginResponse{}, errors.ErrInCredentials
	}

	jwtTTL, err := strconv.Atoi(os.Getenv("JWT_TTL"))
	if err != nil {
		return response.LoginResponse{}, err
	}

	expiresAt := time.Now().Add(time.Hour * time.Duration(jwtTTL)).Unix()

	errf := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	// password does not match
	if errf != nil && errf == bcrypt.ErrMismatchedHashAndPassword {
		return response.LoginResponse{}, errors.ErrInCredentials
	}

	tk := &model.Token{
		UserID:   user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		return response.LoginResponse{}, err
	}
	return response.LoginResponse{User: user, JwtToken: tokenString}, nil
}

// SignUp ...
func (a *AuthApp) SignUp(
	email,
	userName,
	password string,
) (response.SignUpResponse, error) {
	user, err := a.userRepository.GetUserByEmail(email)

	if err != nil {
		return response.SignUpResponse{}, err
	}

	if user.ID > 0 {
		return response.SignUpResponse{}, errors.ErrDupEmail
	}

	jwtTTL, err := strconv.Atoi(os.Getenv("JWT_TTL"))

	if err != nil {
		return response.SignUpResponse{}, err
	}

	expiresAt := time.Now().Add(time.Hour * time.Duration(jwtTTL)).Unix()

	hashPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return response.SignUpResponse{}, err
	}

	hPS := string(hashPass)
	user = model.User{
		Password: hPS,
		Email:    email,
		UserName: userName,
	}

	err = a.userRepository.Add(&user)

	if err != nil {
		return response.SignUpResponse{}, err
	}

	user.Password = ""

	tk := &model.Token{
		UserID:   user.ID,
		UserName: user.UserName,
		Email:    user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	jwtSecret := os.Getenv("JWT_SECRET")

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return response.SignUpResponse{}, nil
	}

	return response.SignUpResponse{User: user, JwtToken: tokenString}, nil

}
