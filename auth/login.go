package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"paypal_clone_project/auth/models"
	"paypal_clone_project/database"
)

func Login(c *gin.Context) {

	var new_user models.User

	if err := c.ShouldBind(&new_user); err != nil {
		c.String(400, "An error has occurred")
		return
	}

	email := new_user.Email       //extracting email from struct model
	password := new_user.Password // extracting password from struct model

	check_login(c, email, password) // check auth credentials

}

func check_login(c *gin.Context, email string, password string) {

	db, err := database.DB_connect()
	if err != nil {
		log.Println(err)
		c.String(500, "An error has occurred")
		return
	}

	var hashed_password string

	err = db.QueryRow("SELECT get_hashed_password($1)", email).Scan(&hashed_password)
	err = bcrypt.CompareHashAndPassword([]byte(hashed_password), []byte(password))
	if err != nil {
		c.String(500, "Password is incorrect.")
		fmt.Println(err)
		return
	}
	var exists bool
	err = db.QueryRow("SELECT check_email_and_password($1, $2)", email, hashed_password).Scan(&exists)
	if err != nil {
		c.String(500, "Error happened")
		fmt.Println(err)
		return
	}

	// return the results
	if exists {
		c.String(200, "Login successful")
		return
	} else {
		c.String(404, "Wrong credentials. Try again")
	}

	err = db.Close()
	if err != nil {
		return
	}

	return
}
