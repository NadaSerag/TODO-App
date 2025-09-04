package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()

	connStr := "postgres://postgres:Mydatabase123@localhost:5432/todo_app?sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	// Testing connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	fmt.Println("Connected to PostgreSQL, Yayy!")
	// router.GET("/", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })

	router.GET("/todos", GetTodos)
	router.GET("/todos/:id", GetTodoById)
	router.POST("/todos", CreateTodo)

	router.Run() // listen and serve on 0.0.0.0:8080
}
