package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/baobei23/todo-realtime-microservices/services/todo-service/internal/domain"
	"github.com/baobei23/todo-realtime-microservices/shared/contracts"
	"github.com/baobei23/todo-realtime-microservices/shared/messaging"
)

// TodoUpdatedEvent is the payload that will be sent to RabbitMQ
type TodoUpdatedEvent struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	UpdatedAt string `json:"updated_at"`
	UserID    string `json:"user_id,omitempty"` // If there is auth in the future
}

type TodoEventPublisher struct {
	amqp *messaging.RabbitMQ
}

func NewTodoEventPublisher(amqp *messaging.RabbitMQ) *TodoEventPublisher {
	return &TodoEventPublisher{
		amqp: amqp,
	}
}

func (p *TodoEventPublisher) PublishTodoUpdated(ctx context.Context, todo *domain.Todo) error {
	// 1. Convert Domain -> Event Payload (Inner Data)
	eventPayload := TodoUpdatedEvent{
		ID:        todo.ID,
		Title:     todo.Title,
		Body:      todo.Body,
		UpdatedAt: todo.UpdatedAt.String(),
	}

	payloadBytes, err := json.Marshal(eventPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// 2. Wrap dalam contracts.AmqpMessage
	msg := contracts.AmqpMessage{
		RoomID: todo.ID, // RoomID can be used to identify the specific room
		Type:   contracts.TodoEventUpdated,
		Data:   payloadBytes,
	}

	log.Printf("Publishing Update Event for Todo %s", todo.ID)

	// 3. Publish menggunakan shared wrapper
	// Note: We don't need to specify exchange, it's handled by the wrapper
	if err := p.amqp.PublishMessage(ctx, contracts.TodoEventUpdated, msg); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}
