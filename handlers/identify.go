package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jipiboily/forwardlytics/integrations"
)

// Identify is taking an identification to send it to the enabled integrations
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
	var identification integrations.Identification
	err := decoder.Decode(&identification)
	if err != nil {
		logrus.WithField("err", err).WithField("body", r.Body).Error("Bad request in Identify")
		writeResponse(w, "Invalid request.", http.StatusBadRequest)
		return
	}
	identification.ReceivedAt = receivedAt

	// Input validation
	missingParameters := identification.Validate()
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
			logrus.Infof("Forwarding idenitify to %s", integrationName)
			err := integration.Identify(identification)
			if err != nil {
				errMsg := fmt.Sprintf("Fatal error during identification with an integration (%s): %s", integrationName, err)
				logrus.WithField("integration", integrationName).WithField("identification", identification).WithField("err", err).Error("Fatal error during identification")
				writeResponse(w, errMsg, 500)
				return
			}
		}

	}

	writeResponse(w, "Forwarding identify to integrations.", http.StatusOK)
}
