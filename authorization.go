package main

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	var newUser User

	//assigning values in the sent JSON to the attrivutes in out newUser struct
	err := c.ShouldBindJSON(&newUser)

	//checking for invalid JSON: missing username or passowrd, or incorrect syntax
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	//have to give it the password as a byte array, so it wraps it with "[]byte..."
	//second parameter is the cost factor,
	//It controls how computationally expensive the hashing will be.
	// Higher cost → slower hash (and slower brute-force attacks).
	// Lower cost → faster but less secure.
	// Default = 10
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)

	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to hash password"})
		return
	}

	//assigning the hashed password to that user
	newUser.Password = string(hash)

	DB.Create(&newUser)

	//returning response only with UserDTO taht doesn't contain the password
	c.JSON(200, ToUserDTO(newUser))
}

func LogIn(c *gin.Context) {
	var loggingUser User

	err := c.ShouldBindJSON(&loggingUser)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	var presentUser User
	result := DB.Where("username = ?", loggingUser.Username).Find(&presentUser)

	if result.RowsAffected == 0 {
		// code 404 for Resource Not Found
		c.JSON(404, gin.H{"error": fmt.Sprintf("No users with username = %s", loggingUser.Username)})
		return
	}

	err2 := bcrypt.CompareHashAndPassword([]byte(presentUser.Password), []byte(loggingUser.Password))

	if err2 != nil {
		c.JSON(401, gin.H{"message": "❌ Unauthorized; Wrong Password"})
		return
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      presentUser.ID,
		"username": presentUser.Username,
		"role":     presentUser.Role,
		//expires in 24 hrs
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err3 := token.SignedString([]byte("super-secret-key"))

	if err3 != nil {
		c.JSON(400, gin.H{"error": "Failed to create token"})
		return
	}

	//returning the JWT (the string token) in the response
	c.JSON(200, gin.H{"token": tokenString})

	//set Cookie ?
}
