package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"paypal_clone_project/auth"
)

func main() {

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current directory:", dir)

	// Run Gin in the correct mode
	gin.SetMode(gin.ReleaseMode)

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default() // Create a new Gin router

	r.POST("/register", auth.Register)

	// Start the server on port 8080
	err = r.Run(":8080")
	if err != nil {
		return
	}
}
