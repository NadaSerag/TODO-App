package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetTodos(t *testing.T) {

	req, _ := http.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()

	// 	router.ServeHTTP(w, req) internally creates the *gin.Context and passes it to GetTodos(c* gin.Context) handler function.
	//  therefore no need to manually create c *gin.Context.
	router.ServeHTTP(w, req)
	// router.ServeHTTP(w, req) checks the registered routes.
	// If /todos is not registered → 404.
	// So it's a MUST to register GetTodos on /todos, the router knows which handler to call.

	//checking if the code is 200 (SUCCESS)
	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	t.Log("Response body:", w.Body.String())

}

// 3ady An Array of STRINGS describing different testcases for the POST method
var testCases = []string{
	"Sucess",
	"Invalid JSON - Missing Fields",
	"Invalid Priority",
	"Past Due Date",
}

func TestCraeteTodo(t *testing.T) {

	//SUCCESS
	t.Run(testCases[0], TestCreateTodo_Success)
	//INVALID JSON - MISSING FIELDS
	t.Run(testCases[1], TestCreateTodo_MissingFields)
	//INVALID PRIORITY
	t.Run(testCases[2], TestCreateTodo_InvalidPriority)
	//PAST DUE DATE
	t.Run(testCases[3], TestCreateTodo_PassedDue)
}

func TestCreateTodo_Success(t *testing.T) {
	//ALTERNATIVE: other way than JSON marshalling,
	//manually writing JSON strings like this for example:
	//   jsonBody := []byte(`{
	//       "title": "New Todo",
	//       "completed": false,
	//       "category": "Work",
	//       "priority": "High"
	//   }`)
	trialToAdd := Todo{Title: "POST request in Go Test!", Completed: &notCompleted, Priority: "High", Category: "Work"}
	body, _ := json.Marshal(trialToAdd)

	// 	bytes.NewBuffer() or strings.NewReader() to create an io.Reader for the request body.
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	// Setting the Content-Type header correctly (application/json) if you use c.ShouldBindJSON(&struct).
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	t.Log("Response body:", w.Body.String())
}

func TestCreateTodo_MissingFields(t *testing.T) {

	trialToAdd := Todo{Title: "Missing Category and Priority fields", Completed: &notCompleted}
	body, _ := json.Marshal(trialToAdd)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
	t.Log("Response body:", w.Body.String())
}

func TestCreateTodo_InvalidPriority(t *testing.T) {

	trialToAdd := Todo{Title: "Invalid priority task", Completed: &notCompleted, Priority: "very HIGHH", Category: "Study"}
	body, _ := json.Marshal(trialToAdd)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
	t.Log("Response body:", w.Body.String())
}

func TestCreateTodo_PassedDue(t *testing.T) {

	due := time.Date(2025, 7, 20, 12, 0, 0, 0, time.UTC)
	trialToAdd := Todo{Title: "Past Due Date", Completed: &notCompleted, Priority: "High", Category: "Study", DueDate: &due}
	body, _ := json.Marshal(trialToAdd)
	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
	t.Log("Response body:", w.Body.String())
}
