package service

import (
	"context"

	"github.com/baobei23/todo-realtime-microservices/services/todo-service/internal/domain"
)

type service struct {
	repo domain.TodoRepository
}

func NewService(repo domain.TodoRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, title, body string) (*domain.Todo, error) {
	todo := &domain.Todo{
		Title: title,
		Body:  body,
	}
	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *service) Get(ctx context.Context, id string) (*domain.Todo, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) List(ctx context.Context, limit, offset int) ([]*domain.Todo, int, error) {
	if limit == 0 {
		limit = 10
	}
	return s.repo.List(ctx, limit, offset)
}

func (s *service) Update(ctx context.Context, id, title, body string) (*domain.Todo, error) {
	existing, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		existing.Title = title
	}
	if body != "" {
		existing.Body = body
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	// Nanti event publish di sini
	return existing, nil
}
