package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testCasesSIGNUP = []string{
	"Successful: ",
	"Invalid JSON: ",
	"Duplicate username: ",
}

func TestSignUp(t *testing.T) {
	//SUCCESSFUL SIGNING UP
	t.Run(testCasesSIGNUP[0], TestSignUp_Success)
	//INVALID JSON
	t.Run(testCasesSIGNUP[1], TestSignUp_Fails)
	//DUPLICATE USERNAME
	t.Run(testCasesSIGNUP[2], TestSignUp_Duplicate)
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

	t.Log("Response body:", w.Body.String())
}

func TestSignUp_Fails(t *testing.T) {
	//INVALID JSON, missing closing quotes for username
	invalidJSON := `{"username": "testingbuddy , "password": "aprettypassword3030"}`
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
	t.Log("Response:", w.Body.String())
}

func TestSignUp_Duplicate(t *testing.T) {
	//INVALID JSON, missing closing quotes for username
	invalidJSON := `{"username": "nadaserag , "password": "aprettypassword3030"}`
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}
	t.Log("Response:", w.Body.String())
}

var testCasesLOGIN = []string{
	"Successful: ",
	"User doesn't have an account: ",
	"Wrong Password: ",
}

func TestLogIn(t *testing.T) {
	//SUCCESSFUL LOGING UP
	t.Run(testCasesLOGIN[0], TestLogIn_Success)
	//USER DIDNT CREATE AN ACCOUNT
	t.Run(testCasesLOGIN[1], TestLogIn_nonExistent)
	//WRONG PASSWORD
	t.Run(testCasesLOGIN[2], TestLogIn_WrongPass)

}
func TestLogIn_Success(t *testing.T) {
	userToSign := User{Username: "nadaaserag", Password: "mypass"}
	body, _ := json.Marshal(userToSign)

	// 	bytes.NewBuffer() or strings.NewReader() to create an io.Reader for the request body.
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	// Setting the Content-Type header correctly (application/json) if you use c.ShouldBindJSON(&struct).
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	t.Log("Token generated:", w.Body.String())
}

func TestLogIn_nonExistent(t *testing.T) {
	userToSign := User{Username: "userNotSignedUP", Password: "lol"}
	body, _ := json.Marshal(userToSign)

	// 	bytes.NewBuffer() or strings.NewReader() to create an io.Reader for the request body.
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	// Setting the Content-Type header correctly (application/json) if you use c.ShouldBindJSON(&struct).
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	t.Log("Token generated:", w.Body.String())
}

func TestLogIn_WrongPass(t *testing.T) {
	userToSign := User{Username: "nadaserag", Password: "wrongpass"}
	body, _ := json.Marshal(userToSign)

	// 	bytes.NewBuffer() or strings.NewReader() to create an io.Reader for the request body.
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	// Setting the Content-Type header correctly (application/json) if you use c.ShouldBindJSON(&struct).
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Fatalf("Expected 200 OK, got %d", w.Code)
	}

	t.Log("Token generated:", w.Body.String())
}
