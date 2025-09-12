package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuthentication(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	//check: If the header is missing OR If it doesn’t start with "Bearer"
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(401, gin.H{"error": "missing or invalid token"})
		c.Abort()
		return
	}

	//extracting the token from the request header which looks like: Authorization: Bearer eyJhbGc6..(token)...
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Parsing and validating token
	user := &User{}
	token, err := jwt.ParseWithClaims(tokenString, user, func(token *jwt.Token) (interface{}, error) {
		return []byte("super-secret-key"), nil
	})

	c.Next()
}

func RequireAuthorization(c *gin.Context) {

	c.Next()
}
