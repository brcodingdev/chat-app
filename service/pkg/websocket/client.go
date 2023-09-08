package websocket

import (
	"encoding/json"
	"github.com/brcodingdev/chat-app/service/pkg/app"
	"log"
	"strings"

	"github.com/gorilla/websocket"
)

// Client ...
type Client struct {
	App        app.Chat
	ID         string
	Connection *websocket.Conn
	Pool       *Pool
	Email      string
	UserID     uint
}

// Message ...
type Message struct {
	Type int `json:"Type,omitempty"`
	Body Body
}

// Body ...
type Body struct {
	ChatRoomName string `json:"chatRoomName,omitempty"`
	ChatRoomId   int32  `json:"chatRoomId,omitempty"`
	ChatMessage  string `json:"chatMessage,omitempty"`
	ChatUser     string `json:"chatUser,omitempty"`
}

// Read ...
func (c *Client) Read(bodyChan chan []byte) {
	defer func() {
		c.Pool.Unregister <- c
		c.Connection.Close()
	}()
	defer c.Pool.ReviveWebsocket()

	for {
		messageType, p, err := c.Connection.ReadMessage()
		if err != nil {
			log.Printf("could not read message %v", err)
			continue
		}

		var body Body
		err = json.Unmarshal(p, &body)
		if err != nil {
			log.Printf("could not unmarshall body %v", err)
			continue
		}

		body.ChatUser = c.Email
		message := Message{Type: messageType, Body: body}
		c.Pool.Broadcast <- message
		log.Println("message received: ", body, "messageType: ", messageType)

		// if stock command, then send message to the broker to be processed
		if strings.Index(body.ChatMessage, "/stock=") == 0 {
			bodyChan <- p
			log.Println("sending message to broker")
		} else {
			// otherwise, add message into database and broadcast via websocket
			go c.App.CreateChatMessage(
				body.ChatMessage,
				uint(body.ChatRoomId),
				c.UserID,
			)
		}
	}
}
