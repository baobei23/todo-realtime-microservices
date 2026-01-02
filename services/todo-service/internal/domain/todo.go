package domain

import (
	"context"
	"time"
)

type Todo struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TodoRepository interface {
	Create(ctx context.Context, todo *Todo) error
	Get(ctx context.Context, id string) (*Todo, error)
	List(ctx context.Context, limit, offset int) ([]*Todo, int, error)
	Update(ctx context.Context, todo *Todo) error
}

type TodoService interface {
	Create(ctx context.Context, title, body string) (*Todo, error)
	Get(ctx context.Context, id string) (*Todo, error)
	List(ctx context.Context, limit, offset int) ([]*Todo, int, error)
	Update(ctx context.Context, id, title, body string) (*Todo, error)
}
