package main

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {

	gin.SetMode(gin.TestMode)
	ConnectToDB()

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
	//router.Run()

	code := m.Run() // <--- runs all tests
	// 	m.Run() executes all tests in the package (TestXXX functions).
	// It returns an exit code:
	// 0 → all tests passed
	// non-zero → some tests failed

	// teardown code before os.Exit()
	sqlDatabase, _ := DB.DB()
	sqlDatabase.Close()

	//os.Exit(code) sends a “pass/fail” signal (depending on the code variable returned from m.Run())to the test runner.
	// The runner reads that code and prints PASS or FAIL in your terminal.
	os.Exit(code)
}

func TestGetTodos(c *gin.Context) {

}
