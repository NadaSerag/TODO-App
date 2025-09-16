package main

import (
	_ "github.com/NadaSerag/TODO-App/docs"
	"github.com/NadaSerag/TODO-App/middleware"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title TODO-App API Documentation w/ Swagger
// @version 1.0
// @Description API documentation of my RESTful API - Todo App that allows clients to manage tasks (todos) through standard HTTP methods.
// @Description Each todo is treated as a resource, identified by a unique ID, and manipulated using predictable endpoints.
// @Description Authentication is handled via **JWT tokens**.
// @Description
// @Description To access protected endpoints, a valid token must be included in the `Authorization` header: `Authorization: Bearer <your_token>`
// @contact.name Nada Serag
// @contact.email nadaserag2006@gmail.com
// @host localhost:8080
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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

	// GetTodos returns list of all todos
	// @Summary      List all todos
	// @Description  Returns a list of all todos
	// @Tags         Todos
	// @Success      200  {array} Todo
	// @Produce json
	// @Router       /todos [get]
	router.GET("/todos", GetTodos)

	// GetTodoById returns a single todo by ID
	// @Summary Get a todo by ID
	// @Description Returns a single todo by its ID
	// @Tags Todos
	// @Param id path int true "Todo ID"
	// @Success 200 {object} Todo
	// @Failure 404 {object} gin.H
	// @Router /todos/{id} [get]
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

	// Swagger UI at /swagger/index.html
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run() // listen and serve on 0.0.0.0:8080

	//closing our database
	sqlDatabase, _ := DB.DB()
	sqlDatabase.Close()

}
