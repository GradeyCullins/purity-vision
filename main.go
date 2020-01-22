package main

import (
	"flag"
	"log"

	"github.com/GradeyCullins/GoogleVisionFilter/src/config"
	"github.com/GradeyCullins/GoogleVisionFilter/src/db"
	"github.com/GradeyCullins/GoogleVisionFilter/src/server"
)

var portFlag int

func main() {
	flag.IntVar(&portFlag, "port", config.DefaultPort, "port to run the service on")
	flag.Parse()

	db, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// rows, err := db.Query("SELECT * from users")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var user models.User
	// 	if err = rows.Scan(&user.UID, &user.Email, &user.Password); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Println(user)
	// }

	server.InitWebServer(portFlag, db)
}
