package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default() // Create a new Gin router

	// Define a simple route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})

	// Start the server on port 8080
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
