package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// the "secret" to generate the JWT with
var JwtSecret = []byte("verylongheybkhdbsuhoeua569u985wcthrq3cjktbx4j")

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
	fmt.Println("tokenString = ", tokenString)
	fmt.Println("authHeader = ", authHeader)

	// Parsing and validating token
	//"claims" is a POINTER to a Claims struct
	claims := &Claims{}

	// 	jwt.ParseWithClaims takes a JWT string verifies it, and decodes the claims into the struct you provide.
	// It does 3 big jobs:
	// Splits the token into header.payload.signature.
	// Verifies the signature using the secret.
	// Decodes the payload into your struct (Claims).
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(401, gin.H{"error": "invalid token"})
		c.Abort()
		return
	}

	//tokenn.ExpiresAt = "1757853813"
	//jwt.NewNumericDate(time.Unix(1757853813, 0)),
	//claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 24))

	// fmt.Println("ExpiresAt:", claims.ExpiresAt) // should NOT be nil
	if claims.ExpiresAt == nil {
		c.JSON(401, gin.H{"error": "expiry date is nil!"})
		c.Abort()
		return
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(401, gin.H{"error": "token expired: 2nd check failure"})
		c.Abort()
		return
	}

	// Putting claims it in context under "user"
	// and Handlers know they can grab "user" to see who is logged in.
	//storing a pointer to Claims in Gin’s context under the key "user".
	c.Set("user", claims)

	fmt.Println(claims)
	fmt.Println("Passed Middleware")
	c.Next()
}

func RequireAuthorization(c *gin.Context) {

	gottenClaims, exists_ok := GetUserClaims(c)

	everythingOk := ClaimsCheck(c, gottenClaims, exists_ok)

	//	if false if returned by ClaimsCheck, then it c.aborted in the function ClaimsCheck already, so we just exit (return)
	if !everythingOk {
		return
	}

	//now we can use the claims returned from GetUserClaims safely without a worry in the world
	if gottenClaims.Role == "admin" {
		fmt.Println("Admin authorized, passed to route successfully")
		//Congrats! Now you're authorized — continue to the route function
		c.Next()
	} else {
		c.JSON(403, gin.H{"message": "Forbidden (user lacks permission to delete ALL todos in the entire database)."})
		c.Abort()
		return
	}
}
