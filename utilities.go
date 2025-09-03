package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type todo struct {
	id        int
	title     string
	completed bool
}

var Todos []todo

func InitializeTodoArray() []todo {
	//The in-memory store (slice) is initialized as empty.
	Todos = []todo{}
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
		if Todos[i].id == id {
			c.JSON(200, Todos[i])
			return
		}
	}

	// code 404 used bec: Todo not found for GET, PUT, or DELETE.
	c.JSON(404, gin.H{"error": "Todo not found"})

}
