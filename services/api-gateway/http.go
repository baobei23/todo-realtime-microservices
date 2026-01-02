package main

import (
	"net/http"
	"strconv"

	pb "github.com/baobei23/todo-realtime-microservices/shared/proto/todo"
	"github.com/gin-gonic/gin"
)

func createTodoHandler(client pb.TodoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Title string `json:"title" binding:"required"`
			Body  string `json:"body"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := client.CreateTodo(c.Request.Context(), &pb.CreateTodoRequest{
			Title: req.Title,
			Body:  req.Body,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, res.Todo)
	}
}

func listTodosHandler(client pb.TodoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

		res, err := client.ListTodos(c.Request.Context(), &pb.ListTodosRequest{
			Limit:  int32(limit),
			Offset: int32(offset),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func getTodoHandler(client pb.TodoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		res, err := client.GetTodo(c.Request.Context(), &pb.GetTodoRequest{
			Id: id,
		})

		if err != nil {
			// Idealnya cek status code grpc (NotFound vs Internal)
			c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found or error occurred"})
			return
		}

		c.JSON(http.StatusOK, res.Todo)
	}
}

func updateTodoHandler(client pb.TodoServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var req struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		res, err := client.UpdateTodo(c.Request.Context(), &pb.UpdateTodoRequest{
			Id:    id,
			Title: req.Title,
			Body:  req.Body,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, res.Todo)
	}
}
