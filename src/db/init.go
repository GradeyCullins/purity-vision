package db

import (
	"context"
	"database/sql"
	"fmt"
	"google-vision-filter/src/config"
	"time"

	"github.com/jackc/pgx/v4"
)

// User represents a user in the Purity system.
type User struct {
	UID      int    `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Image represents a pass/fail status of a previously queried img URI.
type Image struct {
	ImgURIHash string         `json:"imgURIHash"`
	Error      sql.NullString `json:"error"`
	Pass       bool           `json:"pass"`
	DateAdded  time.Time      `json:"dateAdded"`
}

// InitDB intializes and returns a postgres database connection.
func InitDB() (*pgx.Conn, error) {
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
	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}
