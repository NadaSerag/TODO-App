package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// The "jwt.RegisteredClaims" (from github.com/golang-jwt/jwt/v5)
// includes a set of standard JWT fields which are:

// type RegisteredClaims struct {
//     Issuer    string         `json:"iss,omitempty"`
//     Subject   string         `json:"sub,omitempty"`
//     Audience  ClaimStrings   `json:"aud,omitempty"`
//     ExpiresAt *NumericDate   `json:"exp,omitempty"`
//     NotBefore *NumericDate   `json:"nbf,omitempty"`
//     IssuedAt  *NumericDate   `json:"iat,omitempty"`
//     ID        string         `json:"jti,omitempty"`
// }

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
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("super-secret-key"), nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	// Checking expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(401, gin.H{"error": "token expired"})
		c.Abort()
		return
	}

	// Putting claims it in context under "user"
	// and Handlers know they can grab "user" to see who is logged in.
	c.Set("user", claims)
	c.Next()
}

func RequireAuthorization(c *gin.Context) {

	c.Next()
}
