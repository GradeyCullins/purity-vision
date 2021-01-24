package config

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
)

var (
	// DefaultDBName is the default name of the database.
	DefaultDBName = "purity"

	// DefaultDBTestName is the default name of the test database.
	DefaultDBTestName = "purity_test"

	// DefaultPort is the default port to expose the API server.
	DefaultPort int = 8080

	// DBHost is the host machine running the postgres instance.
	DBHost string = getEnvWithDefault("PURITY_DB_HOST", "localhost")

	// DBPort is the port that exposes the db server.
	DBPort string = getEnvWithDefault("PURITY_DB_PORT", "5432")

	// DBName is the postgres database name.
	DBName string = getEnvWithDefault("PURITY_DB_NAME", DefaultDBName)

	// DBUser is the postgres user account.
	DBUser string = getEnvWithDefault("PURITY_DB_USER", "postgres")

	// DBPassword is the password for the DBUser postgres account.
	DBPassword string = getEnvWithDefault("PURITY_DB_PASS", "")

	// DBSSLMode sets the SSL mode of the postgres client.
	DBSSLMode string = getEnvWithDefault("PURITY_DB_SSL_MODE", "disable")

	// LogLevel is the level of logging for the application.
	LogLevel string = getEnvWithDefault("PURITY_LOG_LEVEL", strconv.Itoa(int(zerolog.InfoLevel)))
)

func getEnvWithDefault(name string, def string) string {
	res := os.Getenv(name)
	if res == "" {
		return def
	}
	return res
}
