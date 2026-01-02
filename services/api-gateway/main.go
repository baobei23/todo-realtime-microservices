package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/baobei23/todo-realtime-microservices/services/api-gateway/grpc_clients"
	"github.com/baobei23/todo-realtime-microservices/shared/env"
	"github.com/gin-gonic/gin"
)

var (
	httpAddr = env.GetString("GATEWAY_HTTP_ADDR", ":8081")
)

func main() {
	// 1. Setup gRPC Client menggunakan package baru
	todoClient, err := grpc_clients.NewTodoServiceClient()
	if err != nil {
		log.Fatalf("Failed to init todo grpc client: %v", err)
	}
	defer todoClient.Close()

	r := gin.Default()

	// REST Endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
	v1 := r.Group("/api/v1")
	{
		v1.POST("/todos", createTodoHandler(todoClient.Client))
		v1.GET("/todos", listTodosHandler(todoClient.Client))
		v1.GET("/todos/:id", getTodoHandler(todoClient.Client))
		v1.PUT("/todos/:id", updateTodoHandler(todoClient.Client))
	}

	server := &http.Server{
		Addr:    httpAddr,
		Handler: r,
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("API Gateway listening on %s", httpAddr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		log.Printf("Error starting the API Gateway: %v", err)

	case sig := <-shutdown:
		log.Printf("API Gateway is shutting down due to %v signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Could not stop the API Gateway gracefully: %v", err)
			server.Close()
		}
	}
}
