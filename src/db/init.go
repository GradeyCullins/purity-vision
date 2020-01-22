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
	dbName := config.DB.Name
	dbUser := config.DB.User
	dbPassword := config.DB.Password
	sslMOde := config.DB.SSLMode
	connStr := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=%s", dbName, dbUser, dbPassword, sslMOde)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
