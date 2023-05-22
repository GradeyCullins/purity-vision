package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"purity-vision-filter/src/mail"
	"purity-vision-filter/src/utils"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/webhook"
)

func health(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("All Good ☮️"))
}

func handleFilter(w http.ResponseWriter, req *http.Request) {
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

const MAX_IMAGES_PER_REQUEST = 16

func handleBatchFilter(logger zerolog.Logger) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
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

		// log.Debug().Msgf("Filtering: %v", filterReqPayload.ImgURIList)

		// Validate the request payload URIs
		for _, uri := range filterReqPayload.ImgURIList {
			if _, err := url.ParseRequestURI(uri); err != nil {
				writeError(400, fmt.Sprintf("%s is not a valid URI\n", uri), w)
				return
			}
		}

		for i := 0; i < len(filterReqPayload.ImgURIList); {
			var endIdx int
			if i+MAX_IMAGES_PER_REQUEST > len(filterReqPayload.ImgURIList)-1 {
				endIdx = len(filterReqPayload.ImgURIList)
			} else {
				endIdx = i + MAX_IMAGES_PER_REQUEST
			}

			temp, err := filterImages(
				filterReqPayload.ImgURIList[i:endIdx],
				filterReqPayload.FilterSettings,
			)
			if err != nil {
				logger.Error().Msgf("Error while filtering: %s\n", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			res = append(res, temp...)

			i += MAX_IMAGES_PER_REQUEST
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

type License struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	StripeID string `json:"stripeID"`
	IsValid  bool   `json:"isValid"`
}

func handleWebhook(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Debug().Msgf("Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := "whsec_d42d413b9f320b257032560cf653914f09edc52ee4f479cfaada92ddd402de6c"
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, req.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		logger.Debug().Msgf("Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	// case "charge.succeeded":
	case "invoice.payment_succeeded":
		// Lookup license by customer id
		invoice := stripe.Invoice{}
		err := json.Unmarshal(event.Data.Raw, &invoice)
		if err != nil {
			fmt.Fprintf(w, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		stripeID := invoice.Customer.ID
		email := invoice.CustomerEmail

		license, err := pgStore.GetLicenseByStripeID(stripeID)
		if err != nil {
			logger.Debug().Msgf("Error fetching license: %v\n", err)
			PrintSomethingWrong(w)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// if license exists ensure IsValid is true and return
		if license != nil {
			logger.Debug().Msg("Existing license found, ensuring IsValid is true")
			license.IsValid = true
			if err := pgStore.UpdateLicense(license); err != nil {
				PrintSomethingWrong(w)
				w.WriteHeader(http.StatusInternalServerError)
			}
			// TODO: email person to remind them their subscription is renewed.
			return
		}

		// else create new license and store in db
		logger.Debug().Msg("No license found. Creating one")
		licenseID := utils.GenerateLicenseKey()
		logger.Debug().Msgf("Generated license: %s\n", licenseID)

		license = &License{
			ID:       licenseID,
			Email:    email,
			StripeID: stripeID,
			IsValid:  true,
		}

		if _, err = conn.Model(license).Insert(); err != nil {
			logger.Debug().Msgf("Error creating: %v\n", err)
			PrintSomethingWrong(w)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if err = mail.SendLicenseMail(license.Email, license.ID); err != nil {
			// TODO: retry sending email so user can get their license.
			logger.Debug().Msgf("Error sending license email: %v\n", err)
			PrintSomethingWrong(w)
			w.WriteHeader(http.StatusInternalServerError)
		}
	case "customer.subscription.updated":
		sub := stripe.Subscription{}
		err := json.Unmarshal(event.Data.Raw, &sub)
		if err != nil {
			fmt.Fprintf(w, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		license, err := pgStore.GetLicenseByStripeID(sub.Customer.ID)
		if err != nil {
			logger.Debug().Msgf("Error finding license for valid subscriber: %v\n", err)
		}
		if license == nil {
			logger.Debug().Msg("Failed to find license for existing subscriber. Something is terribly wrong")
			PrintSomethingWrong(w)
			return
		}

		switch sub.Status {
		case stripe.SubscriptionStatusActive:
			license.IsValid = true
			logger.Debug().Msgf("Activated license: %s\n", license.ID)
		case stripe.SubscriptionStatusIncomplete:
		case stripe.SubscriptionStatusIncompleteExpired:
		case stripe.SubscriptionStatusPastDue:
		case stripe.SubscriptionStatusUnpaid:
		case stripe.SubscriptionStatusCanceled:
			license.IsValid = false
			logger.Debug().Msgf("Invalidated license: %s\n", license.ID)
		}

		if err = pgStore.UpdateLicense(license); err != nil {
			logger.Debug().Msgf("Error updating license: %v\n", err)
			PrintSomethingWrong(w)
			return
		}
	default:
		// fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func getLicense(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	licenseID := vars["id"]

	logger.Debug().Msgf("Verifying license: %s\n", licenseID)

	license, err := pgStore.GetLicenseByID(licenseID)
	if err != nil {
		logger.Debug().Msgf("Verifying license: %s\n", licenseID)
		w.WriteHeader(http.StatusInternalServerError)
		PrintSomethingWrong(w)
		return
	}

	if license == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "License %s was not found", licenseID)
		return
	}

	json.NewEncoder(w).Encode(license)
	w.WriteHeader(http.StatusOK)
}
