package auth

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"paypal_clone_project/auth/models"
	"paypal_clone_project/database"
	"time"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var jwt_key = []byte("denis1234")

func Login(c *gin.Context) {

	var new_user models.User

	if err := c.ShouldBind(&new_user); err != nil {
		c.JSON(400, gin.H{"error": "An error has occurred"})
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
		c.String(401, "Invalid email or password")
		log.Print(err)
		return
	}

	token, err := generateJWT(email)
	if err != nil {
		c.String(500, "Interval server error")
		log.Fatal(err)
		return
	}

	c.SetCookie("jwt", token, 3600, "/", "localhost", false, false)

	// return success response
	c.JSON(200, gin.H{
		"message": "Login succesfull",
		"token":   token,
	})

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)
}

func generateJWT(email string) (string, error) {

	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	// create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the token with the secret key
	signedToken, err := token.SignedString(jwt_key)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return signedToken, nil
}

// Middleware to verify JWT token and extract user claims
func VerifyJWT(c *gin.Context) {
	tokenString, err := c.Cookie("jwt")
	fmt.Println("TOKEN", tokenString)
	if err != nil {
		log.Println("EEROR COOKIE", err)
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method")
		}
		// Return the secret key for validation
		return jwt_key, nil
	})

	if err != nil || !token.Valid {
		log.Println("Token invalid", err)
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// Extract claims from the token
	claims, ok := token.Claims.(*Claims)
	if !ok {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		c.Abort()
		return
	}

	// You can now access the email from the claims
	c.Set("email", claims.Email)
	//c.String(200, "token valid")
}

func Logout(c *gin.Context) {
	// clearing jwt cookie by setting it to an expired date
	c.SetCookie("jwt", "", -1, "/", "localhost", false, true)

	c.String(200, "Logout Succesfull")
}

// other GET functions

func GetUserInfo(c *gin.Context) {
	// Retrieve the email from the context set in VerifyJWT
	email, exists := c.Get("email")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized: No email found"})
		return
	}

	// Type assert to string
	userEmail, ok := email.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Internal server error: Invalid email type"})
		return
	}

	db, err := database.DB_connect()
	if err != nil {
		log.Println("Database connection error:", err)
		c.String(500, "Error connecting to database")
		return
	}
	defer db.Close() // Ensure you close the DB connection

	var userBalance float64
	var userCard int64
	var userName string
	var validTo string
	err = db.QueryRow("SELECT * FROM get_user_details($1)", userEmail).Scan(&userBalance, &userCard, &userName, &userEmail, &validTo)
	if err != nil {
		log.Println("Error retrieving balance:", err)
		c.String(500, "Error retrieving user info")
		return
	}

	c.JSON(200, gin.H{
		"balance":  userBalance,
		"card":     userCard,
		"email":    userEmail,
		"username": userName,
		"valid_to": validTo,
	})
}

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email"` // New email to change to
}

func ChangeEmail(c *gin.Context) {

	var request ChangeEmailRequest

	// Bind the JSON payload to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Retrieve the email from the JWT token
	tokenString, err := c.Cookie("jwt")
	if err != nil {
		log.Println("Error obtaining cookie:", err)
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwt_key, nil // Replace with your secret key
	})

	if err != nil || !token.Valid {
		log.Println("Invalid Token:", err)
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	currentEmail := claims.Email // Extract the current email from claims

	// Connect to the database
	db, err := database.DB_connect()
	if err != nil {
		log.Println("Database connection error:", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	defer db.Close() // Ensure the database connection is closed

	// Call the stored procedure to update email
	_, err = db.Exec("SELECT * FROM update_user_email($1, $2)", currentEmail, request.NewEmail)
	if err != nil {
		log.Println("Error updating email:", err)
		c.JSON(500, gin.H{"error": "Failed to change email"})
		return
	}

	c.JSON(200, gin.H{"message": "Email updated successfully"})
}

type TransferFundsRequest struct {
	UserEmail string  `json:"user_email" binding:"required,email"` // User's email (must be valid)
	Amount    float64 `json:"amount" binding:"required,min=0"`     // Amount to transfer (must be positive)
}

func TransferFunds(c *gin.Context) {
	var request TransferFundsRequest

	// Bind the JSON payload to the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	db, err := database.DB_connect()
	if err != nil {
		log.Println("Database connection error:", err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	defer db.Close()

	// Call the stored procedure to transfer funds
	_, err = db.Exec("SELECT transfer_funds($1, $2)", request.UserEmail, request.Amount)
	if err != nil {
		log.Println("Error transferring funds:", err)
		c.JSON(500, gin.H{"error": "Failed to transfer funds"})
		return
	}

	c.JSON(200, gin.H{"message": "Transfer successful"})
}
