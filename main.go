package main

import (
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

	DB.Exec(query)

	//OR :
	//DB.AutoMigrate(&Todo{})
}

func main() {
	router := gin.Default()

	//calling this function to establish connection with our database
	ConnectToDB()

	//crating the `todos` table in the PostgreSQL database on startup if it doesn’t exist.
	initializeTable()

	//GORM's AutoMigrate creates or updates the database schema based on my Go structs.
	DB.AutoMigrate(&User{})

	//in GORM:
	//GET -> Find()
	//POST -> Create()
	//PUT -> Update()
	//DELETE -> Delete()

	router.GET("/todos", GetTodos)
	router.GET("/todos/:id", GetTodoById)
	router.GET("/todos/category/:category", GetTodosByCategory)
	router.GET("/todos/status/:status", GetTodosByStatus)
	router.GET("/todos/search/", GetTodosBySearch)
	router.POST("/todos", CreateTodo)
	router.POST("/signup", SignUp)
	router.POST("/login", LogIn)
	router.PUT("/todos/:id", UpdateTodo)
	router.PUT("/todos/category/:category", UpdateTodosByCategory)
	router.DELETE("/todos/:id", DeleteById)
	router.DELETE("/todos", DeleteAll)

	router.Run() // listen and serve on 0.0.0.0:8080
}
