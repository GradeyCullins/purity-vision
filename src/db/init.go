package db

import (
	"database/sql"
	"fmt"

	"github.com/GradeyCullins/GoogleVisionFilter/src/config"

	// Postgres SQL driver.
	_ "github.com/lib/pq"
)

// User represents a user in the Purity system.
type User struct {
	UID      int    `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ImageCacheEntry represents a pass/fail status of a previously queried img URI.
type ImageCacheEntry struct {
	ImgURIHash string         `json:"imgURIHash"`
	Error      sql.NullString `json:"error"`
	Pass       bool           `json:"pass"`
}

// InitDB intializes and returns a postgres database connection.
func InitDB() (*sql.DB, error) {
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbName := config.DBName
	dbUser := config.DBUser
	dbPassword := config.DBPassword
	if dbPassword == "" {
		return nil, fmt.Errorf("Missing postgres password. Export \"PURITY_DB_PASS=<your_password>\"")
	}
	sslMode := config.DBSSLMode
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPassword, sslMode)
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
