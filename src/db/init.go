package db

import (
	"context"
	"database/sql"
	"fmt"
	"google-vision-filter/src/config"
	"google-vision-filter/src/utils"
	"time"

	"github.com/go-pg/pg/v10"
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

// NewImage returns a new Image instance (wow).
func NewImage(imgURI string, err string, pass bool, dateAdded time.Time) Image {
	return Image{
		ImgURIHash: utils.Hash(imgURI),
		Error:      sql.NullString{String: err, Valid: true},
		Pass:       pass,
		DateAdded:  dateAdded,
	}
}

// InitDB intializes and returns a postgres database connection.
func InitDB(dbName string) (*pg.DB, error) {
	dbHost := config.DBHost
	dbPort := config.DBPort
	dbAddr := fmt.Sprintf("%s:%s", dbHost, dbPort)
	if dbName == "" {
		dbName = config.DBName
	}
	dbUser := config.DBUser
	dbPassword := config.DBPassword
	if dbPassword == "" {
		return nil, fmt.Errorf("Missing postgres password. Export \"PURITY_DB_PASS=<your_password>\"")
	}

	// TODO: use
	// tlsConfig := &tls.Config{}

	conn := pg.Connect(&pg.Options{
		Addr:     dbAddr,
		User:     dbUser,
		Password: dbPassword,
		Database: dbName,
	})

	err := conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
