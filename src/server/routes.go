package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stripe/stripe-go/webhook"
)

// Init intializes the Serve instance and exposes it based on the port parameter.
func (s *Serve) Init(port int, _conn *pg.DB) {
	// Store the database connection in a global var.
	conn = _conn

	// Define handlers.
	batchFilterHandler := http.HandlerFunc(batchFilterImpl)
	filterHandler := http.HandlerFunc(filterImpl)
	healthHandler := http.HandlerFunc(health)
	webhookHandler := http.HandlerFunc(handleWebhook)

	// Create a multiplexer.
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./")))
	mux.Handle("/filter/batch", addCorsHeaders(batchFilterHandler))
	mux.Handle("/filter/single", addCorsHeaders(filterHandler))
	mux.Handle("/health", addCorsHeaders(healthHandler))
	mux.Handle("/webhook", addCorsHeaders(webhookHandler))

	listenAddr = fmt.Sprintf("%s:%d", listenAddr, port)
	log.Info().Msgf("Web server now listening on %s", listenAddr)
	log.Fatal().Msg(http.ListenAndServe(listenAddr, mux).Error())
}

type License struct {
	Id       string `json:"id"`
	StripeID string `json:"stripe_id"`
	IsValid  bool   `json:"is_valid"`
}

func handleWebhook(w http.ResponseWriter, req *http.Request) {
	const MaxBodyBytes = int64(65536)
	req.Body = http.MaxBytesReader(w, req.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
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
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "charge.succeeded":
		fmt.Println("We are making cash")
		licenseID := GenerateLicenseKey()
		fmt.Printf("Generated license: %s\n", licenseID)

		license := License{
			Id:       licenseID,
			StripeID: "some stripe id",
			IsValid:  true,
		}

		_, err := conn.Model(&license).Insert()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("added")
		}
		// Then define and call a function to handle the event payment_intent.succeeded
	// ... handle other event types
	default:
		// fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func GenerateLicenseKey() string {
	key := uuid.New()
	return key.String()
}
