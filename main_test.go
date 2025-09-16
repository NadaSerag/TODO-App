package main

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

var router = gin.Default()

var completed = true
var notCompleted = false

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

	Todo1 := Todo{Title: "Join Meeting", Completed: &completed, Category: "Work", Priority: "High"}
	Todo2 := Todo{Title: "Go Training", Completed: &notCompleted, Category: "Personal", Priority: "Medium"}
	DB.Table("test_todos").Create(&Todo1)
	DB.Table("test_todos").Create(&Todo2)

	//router := gin.Default()

	// router.ServeHTTP(w, req) checks the registered routes.
	// If /todos is not registered → 404.
	// So it's a MUST to register GetTodos on /todos, the router knows which handler to call.
	router.GET("/todos", GetTodos)

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
