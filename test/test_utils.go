package test

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	// a default user variable for the next requests parameter
	user = User{
		UserName: "Valiben",
		Password: "123456",
		Age:      22,
	}

	// a default token for authorization
	myToken = "ssoiuoiu"

	tokenName = "x-xq5-jwt"
)

// receive the ordinary response
type OrdinaryResponse struct {
	Errno  string `json:"errno"`
	Errmsg string `json:"errmsg"`
}

// a middleware function for easy authorization
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(tokenName)
		if token != myToken {
			// unauthorized, just return the unauthorized response
			c.JSON(http.StatusUnauthorized, gin.H{"errno": "-1", "errmsg": "unauthorized"})
			return
		} else {
			// authorized, come into the next function
			c.Next()
		}
	}
}
