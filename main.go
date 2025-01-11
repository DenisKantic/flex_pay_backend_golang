package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"paypal_clone_project/auth"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.POST("/register", auth.Register)
	r.POST("/login", auth.Login)
	r.GET("/logout", auth.Logout)
	r.GET("/token-verify", auth.VerifyJWT)
	r.GET("/protected", auth.VerifyJWT, func(c *gin.Context) {
		email, _ := c.Get("email")
		c.JSON(200, gin.H{
			"Message": "Welcome to the protected route",
			"email":   email,
		})
	})

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
