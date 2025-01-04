package auth

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"paypal_clone_project/database"
	"regexp"
	"time"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CardNum  int    `json:"card_num"`
	ValidTo  string `json:"valid_to"`
}

func Register(c *gin.Context) {

	var newUser User
	if err := c.ShouldBind(&newUser); err != nil {
		c.String(400, "An error has occurred")
		return
	}

	if !is_valid_email(newUser.Email) {
		c.String(400, "Invalid email")
		return
	}

	if len(newUser.Password) < 8 {
		c.String(400, "Password must be at least 8 characters")
		return
	}

	if newUser.CardNum <= 15 {
		c.String(400, "Card number needs to be 16 digits")
		return
	}

	if !check_date(c, newUser.ValidTo) {
		return
	}

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

func is_valid_email(email string) bool {
	regex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func check_date(c *gin.Context, valid_to string) bool {
	today := time.Now()

	fmt.Print("TODAYS DATE", today)
	parsed_date, err := time.Parse("02-01-2006", valid_to)

	fmt.Println("PARSED DATE", parsed_date)
	if err != nil {
		c.String(400, "Invalid date format")
		return false
	}

	if today.After(parsed_date) {
		c.String(400, "You can't choose the older date")
	}

	return true
}
