package middleware

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
)

// EnsureUser zzz
func EnsureUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		userEmail, err := c.Cookie("user-email")

		// Cookie not set so go through the login flow
		if err != nil {
			log.Println("...NO USER-EMAIL - redirect to /login...")
			log.Printf("c.FullPath is: %v\n", c.FullPath())
			log.Printf("c.Request.URL is: %v\n", c.Request.URL)
			// Pass this URI as a query param to the /login endpoint for return redirection
			loginURL, _ := url.Parse("/login")
			loginURLParams := loginURL.Query()
			loginURLParams.Set("redirect_uri", c.Request.URL.String())
			loginURL.RawQuery = loginURLParams.Encode()
			c.Redirect(303, loginURL.String())
			c.Abort()
		} else {
			// Else cookie is set - so add details to Context and continue
			c.Set("user-email", userEmail)
			log.Printf("...user-email = %v...\n", userEmail)
			c.Next()
		}
	}
}
