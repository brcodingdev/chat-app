package broker

import (
	"context"
	"encoding/json"
	"github.com/brcodingdev/chat-app/service/pkg/broker/event"
	"github.com/brcodingdev/chat-app/service/pkg/websocket"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ ...
type RabbitMQ struct {
	receiverQueue  *amqp.Queue
	publisherQueue *amqp.Queue
	channel        *amqp.Channel
	connection     *amqp.Connection
}

// NewRabbitMQ ...
func NewRabbitMQ(
	dsn string,
	receiverQueue string,
	publisherQueue string,
) (*RabbitMQ, error) {

	conn, err := amqp.Dial(dsn)

	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	receiver, err := channel.QueueDeclare(
		receiverQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	publisher, err := channel.QueueDeclare(
		publisherQueue,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &RabbitMQ{
		receiverQueue:  &receiver,
		publisherQueue: &publisher,
		channel:        channel,
		connection:     conn,
	}, nil
}

// Close ...
func (b *RabbitMQ) Close() error {
	if b.channel != nil {
		if err := b.channel.Close(); err != nil {
			return err
		}
	}

	if b.connection != nil {
		if err := b.connection.Close(); err != nil {
			return err
		}
	}

	return nil
}

// Publish send messages to the stock-service
func (b *RabbitMQ) Publish(requestBody chan []byte) {
	for body := range requestBody {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		err := b.channel.PublishWithContext(ctx,
			"",
			b.publisherQueue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			})
		cancel()
		if err != nil {
			log.Printf("could not publish message %s\n", err)
			continue
		}
		log.Printf("message sent %s\n", body)
	}
}

// Consume reads messages from the stock-service
func (b *RabbitMQ) Consume(pool *websocket.Pool) {
	messages, err := b.channel.Consume(
		b.receiverQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("could not read messages %s\n", err)
		return
	}
	receivedStockResponse := make(chan event.StockResponse)
	go messageTransformer(messages, receivedStockResponse)
	go processResponse(receivedStockResponse, pool)
	select {}
}

func messageTransformer(entries <-chan amqp.Delivery, receivedMessages chan event.StockResponse) {
	var stockResponse event.StockResponse
	for entry := range entries {
		if err := json.Unmarshal(entry.Body, &stockResponse); err != nil {
			log.Printf("could not parse response : %s ", string(entry.Body))
			continue
		}

		log.Println("message received")
		receivedMessages <- stockResponse
	}
}

func processResponse(entries <-chan event.StockResponse, pool *websocket.Pool) {
	for entry := range entries {
		log.Println("processing stock response for ", entry.Message)
		stockResponse := event.StockResponse{
			RoomId:  entry.RoomId,
			Message: entry.Message,
		}

		message := websocket.Message{
			Type: 1,
			Body: websocket.Body{ChatRoomId: int32(stockResponse.RoomId),
				ChatUser:    "bot-service",
				ChatMessage: stockResponse.Message,
			},
		}
		// broadcast in the pool channel
		pool.Broadcast <- message
		log.Println("processed", stockResponse.Message)
	}
}
