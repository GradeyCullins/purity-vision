package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// Server listens on localhost:8080 by default.
var listenAddr string = ":8080"

// InitWebServer starts the simple web service.
func InitWebServer() {
	http.HandleFunc("/filter", batchImgFilterHandler)
	log.Printf("Web server now listening on %s\n", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

var batchImgFilterHandler = func(w http.ResponseWriter, req *http.Request) {
	var filterReqPayload batchImgFilterReq

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&filterReqPayload); err != nil {
		writeError(400, "JSON body missing or malformed", w)
		return
	}

	if len(filterReqPayload.ImgURIList) == 0 {
		writeError(400, "ImgUriList cannot be empty", w)
		return
	}

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

	res, err := filterImages(filterReqPayload.ImgURIList)
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
