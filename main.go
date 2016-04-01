package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jipiboily/forwardlytics/integrations"

	_ "github.com/jipiboily/forwardlytics/integrations/drip"
	_ "github.com/jipiboily/forwardlytics/integrations/intercom"
	_ "github.com/jipiboily/forwardlytics/integrations/keen"
	_ "github.com/jipiboily/forwardlytics/integrations/mixpanel"
)

func main() {
	if os.Getenv("FORWARDLYTICS_API_KEY") == "" {
		log.Fatal("You need to set FORWARDLYTICS_API_KEY")
	}

	http.HandleFunc("/identify", identifyHandler)
	log.Println("Forwardlytics started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func identifyHandler(w http.ResponseWriter, r *http.Request) {
	// This is the soonest we can do that, pretty much at least.
	receivedAt := time.Now().Unix()

	// This endpoint is a POST, everything else be a 404
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	// API key validation. Should be moved to a middleware.
	apiKey := r.Header.Get("FORWARDLYTICS_API_KEY")
	if apiKey != os.Getenv("FORWARDLYTICS_API_KEY") {
		log.Printf("Wrong API key. We had '%s' but it should be '%s'\n", apiKey, os.Getenv("FORWARDLYTICS_API_KEY"))

		errorMsg := "Invalid API KEY. The FORWARDLYTICS_API_KEY header must be specified, with the proper API key."
		writeResponse(w, errorMsg, http.StatusUnauthorized)
		return
	}

	// Unmarshal input JSON
	decoder := json.NewDecoder(r.Body)
	var event integrations.Event
	err := decoder.Decode(&event)
	if err != nil {
		log.Println("Bad request:", r.Body)
		writeResponse(w, "Invalid request.", http.StatusBadRequest)
		return
	}
	event.ReceivedAt = receivedAt

	// Input validation
	missingParameters := event.Validate()
	if len(missingParameters) != 0 {
		msg := "Missing parameters: "
		msg = msg + strings.Join(missingParameters, ", ") + "."
		writeResponse(w, msg, http.StatusBadRequest)
		return
	}

	// Yay, it worked so far, let's send all the things to integrations!
	for _, integrationName := range integrations.IntegrationList() {
		integration := integrations.GetIntegration(integrationName)
		if integration.Enabled() {
			log.Println("Forwarding idenitify to", integrationName)
			integration.Identify(event)
		}

	}

	writeResponse(w, "Forwarding identify to integrations.", http.StatusOK)
}

func writeResponse(w http.ResponseWriter, body string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	body = fmt.Sprintf(`{"message": "%s"}`, body)
	w.Write([]byte(body))

}
