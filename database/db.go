package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Importing the pq driver
)

const DB_CONNECT = "postgresql://flex:password@localhost:5434/flexdb?sslmode=disable"

func DB_connect() (*sql.DB, error) {
	db, err := sql.Open("postgres", DB_CONNECT)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	// test the connection

	if err = db.Ping(); err != nil {
		fmt.Println("Error connecting to database:", err)
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	fmt.Println("Successfully connected to database")
	return db, nil
}
