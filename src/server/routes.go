package server

import (
	"fmt"
	"net/http"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

// Init intializes the Serve instance and exposes it based on the port parameter.
func (s *Serve) Init(port int, _conn *pg.DB) {
	// Store the database connection in a global var.
	conn = _conn

	// Define handlers.
	batchFilterHandler := http.HandlerFunc(batchFilterImpl)
	filterHandler := http.HandlerFunc(filterImpl)
	healthHandler := http.HandlerFunc(health)

	// Create a multiplexer.
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world"))
	})
	mux.Handle("/filter/batch", addCorsHeaders(batchFilterHandler))
	mux.Handle("/filter/single", addCorsHeaders(filterHandler))
	mux.Handle("/health", addCorsHeaders(healthHandler))

	listenAddr = fmt.Sprintf("%s:%d", listenAddr, port)
	log.Info().Msgf("Web server now listening on %s", listenAddr)
	log.Fatal().Msg(http.ListenAndServe(listenAddr, mux).Error())
}
