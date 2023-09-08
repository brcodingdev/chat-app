package route

import (
	"encoding/json"
	"fmt"
	"github.com/brcodingdev/chat-app/service/pkg/app"
	"github.com/brcodingdev/chat-app/service/pkg/broker"
	"github.com/brcodingdev/chat-app/service/pkg/websocket"
	"github.com/brcodingdev/chat-app/service/port/http/response"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

// WebSocketProvide ...
type WebSocketProvide struct {
	router *mux.Router
	app    app.Chat
	broker broker.Broker
}

// NewWebSocketProvide ...
func NewWebSocketProvide(
	router *mux.Router,
	app app.Chat,
	broker broker.Broker,
) *WebSocketProvide {
	return &WebSocketProvide{
		router: router,
		app:    app,
		broker: broker,
	}
}

// RegisterRoutes ...
func (c *WebSocketProvide) RegisterRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	c.router.HandleFunc("/ws",
		func(w http.ResponseWriter, r *http.Request) {
			jwtToken := r.URL.Query().Get("jwt")
			jwtSecret := os.Getenv("JWT_SECRET")
			token, err := jwt.Parse(
				jwtToken,
				func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf(
							"unexpected signing method: %v",
							token.Header["alg"],
						)
					}
					return []byte(jwtSecret), nil
				},
			)

			if err != nil {
				handleWebsocketAuthenticationErr(w, err)
				return
			}
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				handleWebsocketAuthenticationErr(w, err)
				return
			}

			c.serveWS(pool, w, r, claims)
		})

}

func (c *WebSocketProvide) serveWS(
	pool *websocket.Pool,
	w http.ResponseWriter,
	r *http.Request,
	claims jwt.MapClaims,
) {
	conn, err := websocket.Upgrade(w, r)

	if err != nil {
		log.Println("could not upgrade websocket")
		return
	}

	client := &websocket.Client{
		App:        c.app,
		Connection: conn,
		Pool:       pool,
		Email:      claims["Email"].(string),
		UserID:     uint(claims["UserID"].(float64)),
	}

	pool.Register <- client
	requestBody := make(chan []byte)
	go client.Read(requestBody)
	go c.broker.Consume(pool)
	go c.broker.Publish(requestBody)
}

func handleWebsocketAuthenticationErr(w http.ResponseWriter, err error) {
	log.Println("websocket error: ", err)
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
