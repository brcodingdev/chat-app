package broker

import (
	"github.com/brcodingdev/chat-app/service/pkg/websocket"
)

// Broker contract to handle pub/sub messages
// In this case we are using RabbitMQ, but we can use other like Kafka
type Broker interface {
	// Publish publishes events
	Publish(requestBody chan []byte)
	// Consume consumes events
	Consume(pool *websocket.Pool)
	// Close to close connections
	Close() error
}
