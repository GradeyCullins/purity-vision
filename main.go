package main

import (
	"flag"
	"google-vision-filter/src/config"
	"google-vision-filter/src/db"
	"google-vision-filter/src/server"
	"log"
)

var portFlag int

func main() {
	flag.IntVar(&portFlag, "port", config.DefaultPort, "port to run the service on")
	flag.Parse()

	conn, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	s := server.NewServe()
	s.Init(portFlag, conn)
}
