package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codeship/go-retro"
	"github.com/jipiboily/forwardlytics/integrations"
)

// Track is taking an event to send it to the enabled integrations
func Track(w http.ResponseWriter, r *http.Request) {
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
		logrus.WithField("err", err).WithField("body", r.Body).Error("Bad request in Track")
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
			logrus.Infof("Forwarding event to %s", integrationName)
			err := retro.DoWithRetry(func() error {
				e := integration.Track(event)
				if e != nil {
					return resourceNotReady(e)
				}
				return e
			})
			if err != nil {
				errMsg := fmt.Sprintf("Fatal error during event with an integration (%s): %s", integrationName, err)
				logrus.WithField("integration", integrationName).WithField("event", event).WithField("err", err).Error("Fatal error during event")
				writeResponse(w, errMsg, 500)
				return
			}
		}
	}

	writeResponse(w, "Forwarding event to integrations.", http.StatusOK)
}
