package server

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

var PrintSomethingWrong = func(w http.ResponseWriter) { fmt.Fprint(w, "Something went wrong") }

// Init intializes the Serve instance and exposes it based on the port parameter.
func (s *Serve) Init(port int, _conn *pg.DB) {
	// Store the database connection in a global var.
	conn = _conn
	pgStore = NewPGStore(conn)

	r := mux.NewRouter()

	r.Use(addCorsHeaders)
	r.Handle("/", http.FileServer(http.Dir("./"))).Methods("GET")
	r.HandleFunc("/health", health).Methods("GET", "OPTIONS")
	r.HandleFunc("/license/{id}", getLicense).Methods("GET")
	r.HandleFunc("/webhook", handleWebhook).Methods("POST")

	// Paywalled filter routes.
	filterR := r.PathPrefix("/filter").Subrouter()
	filterR.Use(paywallMiddleware(pgStore))
	filterR.HandleFunc("/batch", handleBatchFilter(logger)).Methods("POST", "OPTIONS")
	filterR.HandleFunc("/single", handleFilter).Methods("POST", "OPTIONS")

	listenAddr = fmt.Sprintf("%s:%d", listenAddr, port)
	log.Info().Msgf("Web server now listening on %s", listenAddr)
	log.Fatal().Msg(http.ListenAndServe(listenAddr, r).Error())
}
