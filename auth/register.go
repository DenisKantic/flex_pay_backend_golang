package auth

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/smtp"
	"os"
	"paypal_clone_project/auth/models"
	"paypal_clone_project/database"
	"regexp"
	"time"
)

func Register(c *gin.Context) {

	var newUser models.User
	if err := c.ShouldBind(&newUser); err != nil {
		c.String(400, "An error has occurred")
		fmt.Println(err)
		return
	}

	available, message := check_email_name_availability(newUser.Email, newUser.Name)

	if !available {
		c.String(400, message)
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

	if len(newUser.CardNum) != 16 {
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

	// Parse the date string into a time.Time object
	parsedDate, err := time.Parse("02-01-2006", newUser.ValidTo)
	if err != nil {
		log.Fatal("Error parsing date:", err)
	}

	//// create a new random number generator with a custom seed value
	//source := rand.NewSource(time.Now().UnixNano()) // using current unix timestamp as seed
	//r := rand.New(source)
	//
	//activation_code := r.Intn(900000) + 10000
	//fmt.Println("RANDOM NUMBER IS", activation_code)
	//
	//send_activation_code(c, activation_code, newUser.Email) // send activation link to provided email
	//
	////hashing password before putting into register user function

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		c.String(400, "Error has happened. Please contact the admin of the page.")
		fmt.Println(err)
		return
	}

	fmt.Println("hashed password", string(hashed_password))

	var is_profile_activated = true

	_, err = db.Exec("CALL register_user($1,$2,$3,$4,$5,$6)", newUser.Name, newUser.Email,
		hashed_password, newUser.CardNum, parsedDate, is_profile_activated)

	if err != nil {
		c.String(500, "An error has occurred while registering")
		log.Fatal(err)
		return
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}(db)

	send_activation_code(newUser.Email)
	c.JSON(200, gin.H{"success": "Successfully created account"})

}

func is_valid_email(email string) bool {
	regex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func check_date(c *gin.Context, valid_to string) bool {
	today := time.Now()

	fmt.Println("DATE I GOT", valid_to)

	fmt.Print("TODAYS DATE", today)
	parsed_date, err := time.Parse("02-01-2006", valid_to)

	fmt.Println("PARSED DATE", parsed_date)
	if err != nil {
		c.String(400, "Invalid date format")
		return false
	}

	if today.After(parsed_date) {
		c.String(400, "You can't choose the older date")
		return false
	}

	return true
}

func check_email_name_availability(email string, name string) (bool, string) {

	var email_exist, name_exist bool
	db, err := database.DB_connect()

	if err != nil {
		log.Fatal(err)
		return false, "Error has happened."
	}
	err = db.QueryRow("SELECT email_exists, name_exists FROM check_existing_account($1,$2)", email, name).Scan(&email_exist, &name_exist)

	if err != nil {
		fmt.Println(err)
		return false, "Error for registering verification."
	}

	if email_exist {
		return false, "Email already exists"
	}

	if name_exist {
		return false, "Name already exists"
	}

	err = db.Close()
	if err != nil {
		fmt.Println(err)
		return false, "Error has happened."
	}

	return true, ""
}

func send_activation_code(email string) {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	smtp_host := os.Getenv("SMTP_HOST")
	smtp_port := os.Getenv("SMTP_PORT")
	smtp_password := os.Getenv("SMTP_PASSWORD")

	fmt.Println("SMTP_HOST", smtp_host)
	fmt.Println("SMTP_PORT", smtp_port)
	fmt.Println("SMTP_PASSWORd", smtp_password)

	// recipient email
	from := os.Getenv("SMTP_USER")
	to := []string{email}

	// define email message
	subject := "Thank you for choosing our services"
	body := fmt.Sprintf("Welcome to our FlexPay and thank you for choosing us! \n" +
		"\nFlexPay: Simple and easiest way to send and receive money. \n")

	// construct the email
	message := []byte(fmt.Sprintf("To: %s\r\nFrom: %s\r\nSubject: %s\r\n\r\n%s", to[0], from, subject, body))

	// smtp auth
	auth := smtp.PlainAuth("", from, smtp_password, smtp_host)

	// Send the email
	err := smtp.SendMail(smtp_host+":"+smtp_port, auth, from, to, message)
	if err != nil {
		log.Fatal(err)
	}
}
