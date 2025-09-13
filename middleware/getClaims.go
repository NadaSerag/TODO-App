package middleware

import "github.com/gin-gonic/gin"

func GetUserClaims(c *gin.Context) (*Claims, bool) {

	retrievedClaims, exists := c.Get("user")

	if !exists {
		return nil, false
	}

	//since c.Get("...") always returns an interface{} (type 'any'),
	// we need to type-assert it to the type that's originally stored with our c.Set in RequireAuthentication
	claimsStruct, ok := retrievedClaims.(*Claims) // type assertion

	return claimsStruct, ok
}

func ClaimsCheck(c *gin.Context, claimsToBeChecked *Claims, sentBool bool) bool {

	if !sentBool {
		c.JSON(401, gin.H{"error": "user not found in context or claims found in context are of unexpected type"})
		//c.Abort = Don’t continue to the next handlers in the chain after this one.
		c.Abort()
		return false
	}

	//yes, the claims sent to the function are valid, then I can start to use them in the handler function (e.g. RequireAuthorization to check the user's role ...etc)
	return true

}
