package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jipiboily/forwardlytics/integrations"
)

func Identify(w http.ResponseWriter, r *http.Request) {
	// This is the soonest we can do that, pretty much at least.
	receivedAt := time.Now().Unix()

	// This endpoint is a POST, everything else be a 404
	if r.Method != "POST" {
		http.NotFound(w, r)
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
			err := integration.Identify(event)
			if err != nil {
				errMsg := fmt.Sprintf("Fatal error during identification with an integration (%s): %s", integrationName, err)
				log.Println(errMsg)
				writeResponse(w, errMsg, 500)
				return
			}
		}

	}

	writeResponse(w, "Forwarding identify to integrations.", http.StatusOK)
}
