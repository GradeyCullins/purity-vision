package db

import (
	"context"
	"fmt"
	"google-vision-filter/src/config"

	"github.com/go-pg/pg/v10"
)

// User represents a user in the Purity system.
type User struct {
	UID      int    `json:"uid"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Init intializes and returns a postgres database connection object.
func Init(dbName string) (*pg.DB, error) {
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
