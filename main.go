package main

import "github.com/gin-gonic/gin"

func main() {
	router := gin.Default()

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
