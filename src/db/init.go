package db

import (
	"database/sql"
	"fmt"

	"github.com/GradeyCullins/GoogleVisionFilter/src/config"

	// Postgres SQL driver.
	_ "github.com/lib/pq"
)

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
	fmt.Println(connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
