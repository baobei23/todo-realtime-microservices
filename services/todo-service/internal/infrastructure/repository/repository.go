package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/baobei23/todo-realtime-microservices/services/todo-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// InitSchema MVP
func (r *PostgresRepository) InitSchema(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS todos (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        title TEXT NOT NULL,
        body TEXT,
        created_at TIMESTAMP
    );`
	_, err := r.db.Exec(ctx, query)
	return err
}

func (r *PostgresRepository) Create(ctx context.Context, todo *domain.Todo) error {
	query := `
        INSERT INTO todos (title, body, created_at, updated_at) 
        VALUES ($1, $2, $3, $4) 
        RETURNING id`

	todo.CreatedAt = time.Now()

	err := r.db.QueryRow(ctx, query, todo.Title, todo.Body, todo.CreatedAt, todo.UpdatedAt).Scan(&todo.ID)
	if err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (*domain.Todo, error) {
	todo := &domain.Todo{}
	query := `SELECT id, title, body, created_at, updated_at FROM todos WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&todo.ID, &todo.Title, &todo.Body, &todo.CreatedAt, &todo.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("todo not found")
		}
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	return todo, nil
}

func (r *PostgresRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `
        UPDATE todos 
        SET title = $1, body = $2
        WHERE id = $3`

	tag, err := r.db.Exec(ctx, query, todo.Title, todo.Body, todo.ID)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("todo not found")
	}

	return nil
}

func (r *PostgresRepository) List(ctx context.Context, limit, offset int) ([]*domain.Todo, int, error) {
	// 1. Get Total Count
	var total int
	if err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM todos").Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count todos: %w", err)
	}

	// 2. Get Data
	query := `
        SELECT id, title, body, created_at, updated_at 
        FROM todos 
        ORDER BY updated_at DESC 
        LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list todos: %w", err)
	}
	defer rows.Close()

	var todos []*domain.Todo
	for rows.Next() {
		t := &domain.Todo{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Body, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return todos, total, nil
}
