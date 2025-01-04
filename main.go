package main

import (
	"github.com/gin-gonic/gin"
	"paypal_clone_project/auth"
)

func main() {
	r := gin.Default() // Create a new Gin router

	r.POST("/register", auth.Register)

	// Start the server on port 8080
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
