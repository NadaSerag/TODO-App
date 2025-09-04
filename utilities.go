package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Completed   bool   `json:"completed"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
	CompletedAt string `json:"completedAt"`
	DueDate     string `json:"dueDate"`
}

var Todos []Todo
var id = 1

func InitializeTodoArray() []Todo {
	//The in-memory store (slice) is initialized as empty.
	Todos = []Todo{}
	return Todos
}

// - 200: Successful GET, POST, PUT, or DELETE.
// - 400: Invalid JSON or empty title.
// - 404: Todo not found for GET, PUT, or DELETE.

// When a request comes in, Gin creates a *gin.Context object
// and passes it into the function as the parameter
// the parameter is by convention named "c".

func GetTodos(c *gin.Context) {
	//200
	c.JSON(200, Todos)
}

func GetTodoById(c *gin.Context) {
	//c.Param(...) returns the value as a string.
	//thats wy we need to convert it to an int by strconv.Atoi
	id, _ := strconv.Atoi(c.Param("id"))

	for i := 0; i < len(Todos); i++ {
		if Todos[i].Id == id {
			c.JSON(200, Todos[i])
			return
		}
	}

	// code 404 used bec: Todo not found for GET, PUT, or DELETE.
	c.JSON(404, gin.H{"error": "Todo not found"})
}

func CreateTodo(c *gin.Context) {
	var newTodo Todo

	//VERY IMPORTANT: reading the request body
	error := c.ShouldBindJSON(&newTodo)

	//checking for invalid JSON
	if error != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	//checking if the title is empty
	if newTodo.Title == "" {
		c.JSON(400, gin.H{"error": "Title cannot be empty"})
		return
	}

	newTodo.Id = id
	Todos = append(Todos, newTodo)
	id++

	c.JSON(200, newTodo)
}
