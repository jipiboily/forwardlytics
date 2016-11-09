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

// Page is taking a pageview to send it to the enabled integrations
func Page(w http.ResponseWriter, r *http.Request) {
	// This is the soonest we can do that, pretty much at least.
	receivedAt := time.Now().Unix()

	// This endpoint is a POST, everything else be a 404
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	// Unmarshal input JSON
	decoder := json.NewDecoder(r.Body)
	var page integrations.Page
	err := decoder.Decode(&page)
	if err != nil {
		logrus.WithField("err", err).WithField("body", r.Body).Error("Bad request in Page")
		writeResponse(w, "Invalid request.", http.StatusBadRequest)
		return
	}
	page.ReceivedAt = receivedAt

	// Input validation
	missingParameters := page.Validate()
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
			logrus.Infof("Forwarding page to %s", integrationName)
			err := retro.DoWithRetry(func() error {
				e := integration.Page(page)
				if e != nil {
					return resourceNotReady(e)
				}
				return e
			})
			if err != nil {
				errMsg := fmt.Sprintf("Fatal error during page with an integration (%s): %s", integrationName, err)
				logrus.WithField("integration", integrationName).WithField("page", page).WithField("err", err).Error("Fatal error during page")
				writeResponse(w, errMsg, 500)
				return
			}
		}
	}

	writeResponse(w, "Forwarding page to integrations.", http.StatusOK)
}
