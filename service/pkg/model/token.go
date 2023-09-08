package model

import "github.com/golang-jwt/jwt"

// Token ...
type Token struct {
	UserID   uint
	UserName string
	Email    string
	*jwt.StandardClaims
}
