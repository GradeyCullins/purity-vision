package main

import (
	"flag"
	"os"
	"purity-vision-filter/src/config"
	"purity-vision-filter/src/db"
	"purity-vision-filter/src/server"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var portFlag int

func main() {
	flag.IntVar(&portFlag, "port", config.DefaultPort, "port to run the service on")
	flag.Parse()

	logLevel, err := strconv.Atoi(config.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	conn, err := db.Init(config.DefaultDBName)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	defer conn.Close()

	s := server.NewServe()
	s.Init(portFlag, conn)
}
