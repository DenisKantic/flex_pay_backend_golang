package auth

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"paypal_clone_project/database"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CardID   int    `json:"card_id"`
}

func Register(c *gin.Context) {

	var newUser User
	if err := c.ShouldBind(&newUser); err != nil {
		c.String(400, "An error has occurred")
		return
	}
	//test

	db, err := database.DB_connect()

	if err != nil {
		log.Fatal(err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
	c.String(200, "Successfully created user")
	c.String(200, "Successfully created account")

}
