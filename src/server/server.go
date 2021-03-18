package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Server listens on localhost:8080 by default.
var listenAddr string = ""

// Store the db connection passed from main.go.
var conn *pg.DB

var logger zerolog.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

// BatchImgFilterReq is the form of an incoming JSON payload
// for retrieving pass/fail status of each supplied image URI.
type BatchImgFilterReq struct {
	ImgURIList     []string `json:"imgURIList"`
	FilterSettings `json:"filterSettings"`
}

// ImgFilterRes returns the pass/fail status and any errors for a single image URI.
//
// TODO: consider adding 'URI' as a key. Currently, we are taking advantage of the fact
// that the JSON RFC guarantees array ordering. Comment in the below link suggests that
// there are known cases where order is not guaranteed.
// https://stackoverflow.com/a/7214312
type ImgFilterRes struct {
	ImgURI string `json:"imgURI"` // The original URI where the image resides on the web.
	Error  string `json:"error"`  // Any error returned when trying to filter the image.
	Pass   bool   `json:"pass"`   // True if the image passed the filter.
	Reason string `json:"reason"` // Why the image did not pass the filter if non-empty e.g. "adult", "violence", etc.
}

// BatchImgFilterRes represents a list of pass/fail statuses and any errors for each
// supplied image URI.
type BatchImgFilterRes []ImgFilterRes

// ErrorRes is a JSON response containing an error message from the API.
type ErrorRes struct {
	Message string `json:"message"`
}

// Server defines the actions of a Purity API Web Server.
type Server interface {
	Init(int, *sql.DB)
}

// Serve is an instance of a Purity API Web Server.
type Serve struct {
}

// NewServe returns an uninitialized Serve instance.
func NewServe() *Serve {
	return &Serve{}
}

// Init intializes the Serve instance and exposes it based on the port parameter.
func (s *Serve) Init(port int, _conn *pg.DB) {
	// Store the database connection in a global var.
	conn = _conn

	// Define handlers.
	batchFilterHandler := http.HandlerFunc(batchFilter)

	// Create a multiplexer.
	mux := http.NewServeMux()
	mux.Handle("/filter", addCorsHeaders(batchFilterHandler))

	listenAddr = fmt.Sprintf("%s:%d", listenAddr, port)
	log.Info().Msgf("Web server now listening on %s", listenAddr)
	log.Fatal().Msg(http.ListenAndServe(listenAddr, mux).Error())
}

var addCorsHeaders = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token"
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", allowedHeaders)
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		next.ServeHTTP(w, r)
	})
}

func batchFilter(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		var filterReqPayload BatchImgFilterReq

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&filterReqPayload); err != nil {
			writeError(400, "JSON body missing or malformed", w)
			return
		}

		if len(filterReqPayload.ImgURIList) == 0 {
			writeError(400, "ImgUriList cannot be empty", w)
			return
		}

		// Validate the request payload URIs
		for _, uri := range filterReqPayload.ImgURIList {
			if _, err := url.ParseRequestURI(uri); err != nil {
				writeError(400, fmt.Sprintf("%s is not a valid URI\n", uri), w)
				return
			}
		}

		res, err := filter(filterReqPayload)
		if err != nil {
			log.Error().Err(err).Msg("")
			writeError(500, "Something went wrong", w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func writeError(code int, message string, w http.ResponseWriter) {
	logger.Info().Msg(message)
	w.WriteHeader(code)
	err := ErrorRes{
		Message: message,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
