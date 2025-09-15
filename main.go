package main

import (
	"github.com/NadaSerag/TODO-App/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	router := gin.Default()

	//calling this function to establish connection with our database
	ConnectToDB()

	//crating the `todos` table in the PostgreSQL database on startup if it doesn’t exist.
	DB.AutoMigrate(&Todo{})

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
	router.POST("/todos", middleware.RequireAuthentication, CreateTodo)
	router.POST("/signup", SignUp)
	router.POST("/login", LogIn)
	router.PUT("/todos/:id", middleware.RequireAuthentication, UpdateTodo)
	router.PUT("/todos/category/:category", middleware.RequireAuthentication, UpdateTodosByCategory)
	router.DELETE("/todos/:id", middleware.RequireAuthentication, DeleteById)
	router.DELETE("/todos", middleware.RequireAuthentication, middleware.RequireAuthorization, DeleteAll)

	router.Run() // listen and serve on 0.0.0.0:8080
}
