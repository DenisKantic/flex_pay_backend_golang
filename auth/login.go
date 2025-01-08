package auth

import (
	"database/sql"
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

	check_login(c, new_user.Email, new_user.Password) // check auth credentials

	return

}

func check_login(c *gin.Context, email string, password string) {
	// Connect to the database
	db, err := database.DB_connect()
	if err != nil {
		log.Println(err)
		c.String(500, "An error has occurred")
		return
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Query to get email and hashed password from the database
	var hashedPassword string
	err = db.QueryRow("SELECT get_user_email_and_password($1)", email).Scan(&hashedPassword)
	if err != nil {
		if err.Error() == "pq: Email not found" {
			log.Println(err)
			c.String(401, "Invalid email or password")
			return
		}
		log.Println("Database query error:", err)
		c.String(500, "Internal server error")
		return
	}

	// Compare the provided password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Passwords don't match
		c.String(401, "Incorrect password")
		log.Fatal(err)
		return
	}

	// Successful login
	c.String(200, "Login successful")

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
}
