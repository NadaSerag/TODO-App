package main

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	//setting up the TEST database (for testing purposes only)
	query := `
    CREATE TABLE IF NOT EXISTS test_todos (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        completed BOOLEAN DEFAULT FALSE,
        category TEXT,
        priority TEXT,
        completedAt TIMESTAMP NULL,
        dueDate TIMESTAMP NULL
    );`

	DB.Exec(query)

	completed := true
	notCompleted := false
	Todo1 := Todo{Title: "Join Meeting", Completed: &completed, Category: "Work", Priority: "High"}
	Todo2 := Todo{Title: "Go Training", Completed: &notCompleted, Category: "Personal", Priority: "Medium"}
	DB.Create(&Todo1)
	DB.Create(&Todo2)

	//router := gin.Default()

	//router.GET("/todos", TestGetTodos)
	//router.GET("/todos/:id", TestGetTodoById)
}
