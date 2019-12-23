package src

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

// InitWebServer starts the simple web service.
func InitWebServer() {
	http.HandleFunc("/filter", batchImgFilterHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var batchImgFilterHandler = func(w http.ResponseWriter, req *http.Request) {
	var filterReqPayload filterReq

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&filterReqPayload); err != nil {
		log.Fatal(err)
	}

	if len(filterReqPayload.ImgURIList) == 0 {
		w.WriteHeader(400)
		w.Write([]byte("ImgUriList cannot be empty"))
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
			w.WriteHeader(400)
			w.Write([]byte(fmt.Sprintf("%s is not a valid URI\n", uri)))
			return
		}
	}

	filterRes, err := filterImages(filterReqPayload.ImgURIList)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filterRes)
}
