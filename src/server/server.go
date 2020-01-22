package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/GradeyCullins/GoogleVisionFilter/src/filter"
	"github.com/GradeyCullins/GoogleVisionFilter/src/types"
)

// Server listens on localhost:8080 by default.
var listenAddr string = "127.0.0.1"

// Store the db connection passed from main.go.
var db *sql.DB

// InitWebServer starts the simple web service.
func InitWebServer(port int, _db *sql.DB) {
	http.HandleFunc("/filter", batchImgFilterHandler)
	listenAddr = fmt.Sprintf("%s:%d", listenAddr, port)

	db = _db

	log.Printf("Web server now listening on %s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

var batchImgFilterHandler = func(w http.ResponseWriter, req *http.Request) {
	var filterReqPayload types.BatchImgFilterReq

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&filterReqPayload); err != nil {
		writeError(400, "JSON body missing or malformed", w)
		return
	}

	if len(filterReqPayload.ImgURIList) == 0 {
		writeError(400, "ImgUriList cannot be empty", w)
		return
	}

	// TODO: There is a max number of URIs that can be included per batch request,
	// add logic to split the list into multiple requests if that limit is reached.
	// TODO: If the Google Vision API doesn't fail elegantly when the URIs don't
	// point to images, do more validation here that determines whether the URI
	// is or is not an image.
	//
	// Use content-type: image/* as a simple check, but this could be spoofed.
	//
	// Validate the request payload URIs
	for _, uri := range filterReqPayload.ImgURIList {
		if _, err := url.ParseRequestURI(uri); err != nil {
			writeError(400, fmt.Sprintf("%s is not a valid URI\n", uri), w)
			return
		}
	}

	res, err := filter.Images(filterReqPayload.ImgURIList)
	if err != nil {
		fmt.Println(err)
		writeError(500, "Something went wrong", w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func writeError(code int, message string, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Write([]byte(message))
}
