package models

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	CardNum  string `json:"card_num"`
	ValidTo  string `json:"valid_to"`
}
