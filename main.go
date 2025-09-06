package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()

	connectionStr := "postgres://postgres:Mydatabase123@localhost:5432/todo_app?sslmode=disable"

	db, err := sql.Open("postgres", connectionStr)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Connected to PostgreSQL, Yayy!")
	// router.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	router.GET("/todos", GetTodos)
	router.GET("/todos/:id", GetTodoById)
	router.GET("/todos/category/:category", GetTodosByCategory)
	router.GET("/todos/status/:status", GetTodosByStatus)
	router.GET("/todos/search/:q", GetTodosBySearch)
	router.POST("/todos", CreateTodo)

	router.Run() // listen and serve on 0.0.0.0:8080

	db.Close()
}
