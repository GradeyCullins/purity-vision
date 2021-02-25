package db

import (
	"context"
	"fmt"
	"os"
	"purity-vision-filter/src/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/go-pg/pg/v10"
)

var logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

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

	// Print SQL queries to logger if loglevel is set to debug.
	conn.AddQueryHook(loggerHook{})

	err := conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type loggerHook struct{}

func (h loggerHook) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if evt.Err != nil {
		log.Debug().Msgf("%s executing a query:\n%s\n", evt.Err, q)
	} else {
		log.Debug().Msg(string(q))
	}

	return ctx, nil
}

func (loggerHook) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}
