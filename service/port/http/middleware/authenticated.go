package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/brcodingdev/chat-app/service/errors"
	"github.com/brcodingdev/chat-app/service/port/api"
	"github.com/brcodingdev/chat-app/service/port/http/response"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Authenticated ...
func Authenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(
			r.Header.Get("Authorization"),
			"Bearer ",
		)
		jwtSecret := os.Getenv("JWT_SECRET")
		if len(authHeader) != 2 {
			handleAuthenticationErr(w, errors.ErrToken)
			return
		}

		jwtToken := authHeader[1]
		token, err := jwt.Parse(
			jwtToken,
			func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil,
						fmt.Errorf(
							"unexpected signing method: %v",
							token.Header["alg"],
						)
				}
				return []byte(jwtSecret), nil
			})

		if err != nil {
			handleAuthenticationErr(w, err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var props api.JWTProps = "JWTProps"
			ctx := context.WithValue(r.Context(), props, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			handleAuthenticationErr(w, err)
			return
		}
	})
}

func handleAuthenticationErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	res := response.ErrorResponse{
		Message: err.Error(),
		Status:  false,
		Code:    http.StatusUnauthorized,
	}
	data, _ := json.Marshal(res)
	w.Write(data)
}
