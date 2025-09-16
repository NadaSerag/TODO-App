package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestCraeteTodo(t *testing.T) {
	t.Run("Success", func(t *testing.T) { /* ... */ })
	t.Run("EmptyDB", func(t *testing.T) { /* ... */ })
}

func TestCreateTodo_Success(t *testing.T) {

	trialToAdd := Todo{Title: "New Todo", Completed: &notCompleted}
	body, _ := json.Marshal(trialToAdd)
	//ALTERNATIVE: other way than JSON marshalling,
	//manually writing JSON strings like this for example:
	//   jsonBody := []byte(`{
	//       "title": "New Todo",
	//       "completed": false,
	//       "category": "Work",
	//       "priority": "High"
	//   }`)

	req, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	t.Log("Response body:", w.Body.String())

	// 	Use bytes.NewBuffer() or strings.NewReader() to create an io.Reader for the request body.
	// Set the Content-Type header correctly (application/json) if you use c.ShouldBindJSON(&struct).
}
