package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

func health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("All Good ☮️"))
}

func filterImpl(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		var filterReq ImgFilterReq

		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&filterReq); err != nil {
			writeError(400, "JSON body missing or malformed", w)
			return
		}

		if _, err := url.ParseRequestURI(filterReq.ImgURI); err != nil {
			writeError(400, fmt.Sprintf("%s is not a valid URI\n", filterReq.ImgURI), w)
			return
		}

		logger.Info().Msg(filterReq.ImgURI)
		res, err := filterSingle(filterReq.ImgURI, filterReq.FilterSettings)
		if err != nil {
			writeError(500, "Something went wrong", w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func batchFilterImpl(w http.ResponseWriter, req *http.Request) {
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

		var res BatchImgFilterRes

		log.Debug().Msgf("Filtering: %v", filterReqPayload.ImgURIList)

		// Validate the request payload URIs
		for _, uri := range filterReqPayload.ImgURIList {
			if _, err := url.ParseRequestURI(uri); err != nil {
				writeError(400, fmt.Sprintf("%s is not a valid URI\n", uri), w)
				return
			}

			singleRes, err := filterSingle(uri, filterReqPayload.FilterSettings)
			if err != nil {
				log.Error().Msgf("Error while filtering %s: %s", uri, err)
			} else {
				res = append(res, *singleRes)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}
