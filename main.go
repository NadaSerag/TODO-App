package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func initializeTable() {
	//'query' parameter variable for db.Exec(...)
	query := `
    CREATE TABLE IF NOT EXISTS todos (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        completed BOOLEAN DEFAULT FALSE,
        category TEXT,
        priority TEXT,
        completedAt TIMESTAMP NULL,
        dueDate TIMESTAMP NULL
    );`

	db.Exec(query)
}

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

	//crating the `todos` table in the PostgreSQL database on startup if it doesn’t exist.
	initializeTable()

	router.GET("/todos", GetTodos)
	router.GET("/todos/:id", GetTodoById)
	router.GET("/todos/category/:category", GetTodosByCategory)
	router.GET("/todos/status/:status", GetTodosByStatus)
	router.GET("/todos/search/", GetTodosBySearch)
	router.POST("/todos", CreateTodo)
	router.PUT("/todos/:id", UpdateTodo)
	router.PUT("/todos/category/:category", UpdateTodosByCategory)
	router.DELETE("/todos/:id", DeleteById)
	router.DELETE("/todos", DeleteAll)

	router.Run() // listen and serve on 0.0.0.0:8080

	db.Close()
}
