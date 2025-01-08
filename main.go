package main

import (
	"github.com/gin-gonic/gin"
	"paypal_clone_project/auth"
)

func main() {
	// Run Gin in the correct mode
	gin.SetMode(gin.ReleaseMode)

	//if err := godotenv.Load(".env"); err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	r := gin.Default() // Create a new Gin router

	r.POST("/register", auth.Register) // register user account
	r.POST("/login", auth.Login)       // login in user account
	r.GET("/protected", auth.VerifyJWT, func(c *gin.Context) {
		email, _ := c.Get("email")
		c.JSON(200, gin.H{
			"Message": "Welcome to the protected route",
			"email":   email,
		})
	})

	// Start the server on port 8080
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
