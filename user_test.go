package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignUp(t *testing.T) {

	//SUCCESSFUL SIGNING UP
	t.Run("Successful", TestSignUp_Success)
	//INVALID JSON
	t.Run("Failure", TestCreateTodo_MissingFields)

}

func TestSignUp_Success(t *testing.T) {
	userToSign := User{Username: "testingbuddy", Password: "aprettypassword3030$"}
	body, _ := json.Marshal(userToSign)

	// 	bytes.NewBuffer() or strings.NewReader() to create an io.Reader for the request body.
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(body))
	// Setting the Content-Type header correctly (application/json) if you use c.ShouldBindJSON(&struct).
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	t.Log("Token generated:", w.Body.String())
}
