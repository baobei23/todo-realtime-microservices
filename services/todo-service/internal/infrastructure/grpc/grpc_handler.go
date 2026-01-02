package grpc

import (
	"context"
	"time"

	"github.com/baobei23/todo-realtime-microservices/services/todo-service/internal/domain"
	pb "github.com/baobei23/todo-realtime-microservices/shared/proto/todo"
	"google.golang.org/grpc"
)

type gRPCHandler struct {
	pb.UnimplementedTodoServiceServer
	service domain.TodoService
}

func NewGRPCHandler(server *grpc.Server, service domain.TodoService) *gRPCHandler {
	handler := &gRPCHandler{service: service}

	pb.RegisterTodoServiceServer(server, handler)
	return handler
}

func (h *gRPCHandler) CreateTodo(ctx context.Context, req *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	todo, err := h.service.Create(ctx, req.Title, req.Body)
	if err != nil {
		return nil, err
	}
	return &pb.CreateTodoResponse{Todo: convertToProto(todo)}, nil
}

func (h *gRPCHandler) GetTodo(ctx context.Context, req *pb.GetTodoRequest) (*pb.GetTodoResponse, error) {
	todo, err := h.service.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetTodoResponse{Todo: convertToProto(todo)}, nil
}

func (h *gRPCHandler) ListTodos(ctx context.Context, req *pb.ListTodosRequest) (*pb.ListTodosResponse, error) {
	todos, total, err := h.service.List(ctx, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	var protoTodos []*pb.Todo
	for _, t := range todos {
		protoTodos = append(protoTodos, convertToProto(t))
	}

	return &pb.ListTodosResponse{
		Todos:      protoTodos,
		TotalCount: int32(total),
	}, nil
}

func (h *gRPCHandler) UpdateTodo(ctx context.Context, req *pb.UpdateTodoRequest) (*pb.UpdateTodoResponse, error) {
	todo, err := h.service.Update(ctx, req.Id, req.Title, req.Body)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateTodoResponse{Todo: convertToProto(todo)}, nil
}

func convertToProto(t *domain.Todo) *pb.Todo {
	return &pb.Todo{
		Id:        t.ID,
		Title:     t.Title,
		Body:      t.Body,
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
	}
}
