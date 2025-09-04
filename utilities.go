package main

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	Id          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Completed   *bool  `json:"completed" binding:"required"`
	Category    string `json:"category" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	CompletedAt string `json:"completedAt"`
	DueDate     string `json:"dueDate" binding:"required"`
}

// type Todo struct {
// 	Id          int    `json:"id"`
// 	Title       string `json:"title" binding:"required"`
// 	Completed   *bool   `json:"completed"`
// 	Category    string `json:"category"`
// 	Priority    string `json:"priority"`
// 	CompletedAt string `json:"completedAt"`
// 	DueDate     string `json:"dueDate"`
// }

var connectionStr = "postgres://postgres:Mydatabase123@localhost:5432/todo_app?sslmode=disable"
var db, _ = sql.Open("postgres", connectionStr)
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
	//200 code for Successful GET
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
	//returns an error, error is = nil if JSON parses successfully
	errorVar := c.ShouldBindJSON(&newTodo)

	//checking for invalid JSON
	if errorVar != nil {
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

	query := `
        INSERT INTO todos (id, title, completed, category, priority, completedAt, dueDate)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	db.QueryRow(query,
		newTodo.Id,
		newTodo.Title,
		newTodo.Completed,
		newTodo.Category,
		newTodo.Priority,
		newTodo.CompletedAt,
		newTodo.DueDate,
	)
	c.JSON(200, newTodo)
}
