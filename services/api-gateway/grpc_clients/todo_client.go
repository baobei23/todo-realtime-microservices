package grpc_clients

import (
	"os"

	pb "github.com/baobei23/todo-realtime-microservices/shared/proto/todo"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type todoServiceClient struct {
	Client pb.TodoServiceClient
	conn   *grpc.ClientConn
}

func NewTodoServiceClient() (*todoServiceClient, error) {
	todoServiceURL := os.Getenv("TODO_SERVICE_URL")
	if todoServiceURL == "" {
		todoServiceURL = "todo-service:8082"
	}

	// dialOptions := append(
	// 	tracing.DialOptionsWithTracing(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )

	conn, err := grpc.NewClient(todoServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewTodoServiceClient(conn)

	return &todoServiceClient{
		Client: client,
		conn:   conn,
	}, nil
}

func (c *todoServiceClient) Close() {
	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			return
		}
	}
}
