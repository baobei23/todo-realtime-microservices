package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/baobei23/todo-realtime-microservices/shared/contracts"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TodoExchange = "todo"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}

	rmq := &RabbitMQ{
		conn:    conn,
		Channel: ch,
	}

	// Setup topology dasar (Exchange)
	if err := rmq.setupTopology(); err != nil {
		rmq.Close()
		return nil, err
	}

	return rmq, nil
}

func (r *RabbitMQ) setupTopology() error {
	return r.Channel.ExchangeDeclare(
		TodoExchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {
	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	return r.Channel.PublishWithContext(ctx,
		TodoExchange, // exchange
		routingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jsonMsg,
		},
	)
}

func (r *RabbitMQ) ConsumeMessages(queueName string, routingKeys []string, handler func(contracts.AmqpMessage) error) error {
	// 1. Declare Queue
	q, err := r.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// 2. Bind Queue to Exchange for each Routing Key
	for _, key := range routingKeys {
		if err := r.Channel.QueueBind(q.Name, key, TodoExchange, false, nil); err != nil {
			return err
		}
	}

	// 3. Consume
	msgs, err := r.Channel.Consume(
		q.Name,
		"",    // consumer tag
		false, // auto-ack (kita akan ack manual jika sukses)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	// 4. Handle Messages in Background
	go func() {
		for d := range msgs {
			var msg contracts.AmqpMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Error unmarshal message: %v", err)
				d.Nack(false, false) // Reject, jangan requeue (poison message)
				continue
			}

			if err := handler(msg); err != nil {
				log.Printf("Error processing message: %v", err)
				// Di production, gunakan retry/dead-letter.
				// Untuk MVP: Nack + Requeue true (coba lagi) atau false (buang)
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	return nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
}
