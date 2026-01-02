package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/todo-realtime-microservices/services/todo-service/internal/infrastructure/events"
	"github.com/baobei23/todo-realtime-microservices/services/todo-service/internal/infrastructure/grpc"
	"github.com/baobei23/todo-realtime-microservices/services/todo-service/internal/infrastructure/repository"
	"github.com/baobei23/todo-realtime-microservices/services/todo-service/internal/service"
	"github.com/baobei23/todo-realtime-microservices/shared/db"
	"github.com/baobei23/todo-realtime-microservices/shared/env"
	"github.com/baobei23/todo-realtime-microservices/shared/messaging"
	grpcserver "google.golang.org/grpc"
)

func main() {
	dbURI := env.GetString("POSTGRES_URI", "postgresql://postgres:postgres@postgres:5432/todo_db")

	log.Println("Connecting to database...")
	pool, err := db.New(dbURI, 10, 5, 10*time.Second, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// 2. Setup Repository
	repo := repository.NewPostgresRepository(pool)

	// Init Schema
	if err := repo.InitSchema(context.Background()); err != nil {
		log.Fatalf("Failed to init schema: %v", err)
	}

	// 3. Setup Service Logic
	svc := service.NewService(repo)

	// 4. Setup gRPC Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// RabbitMQ connection
	rabbitMqURI := env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672/")
	rabbitmq, err := messaging.NewRabbitMQ(rabbitMqURI)
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitmq.Close()
	log.Println("Starting RabbitMQ connection")

	publisher := events.NewTodoEventPublisher(rabbitmq)

	grpcServer := grpcserver.NewServer()
	grpc.NewGRPCHandler(grpcServer, svc, publisher)

	// 5. Graceful Shutdown
	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	grpcServer.GracefulStop()
	log.Println("Server stopped")
}
