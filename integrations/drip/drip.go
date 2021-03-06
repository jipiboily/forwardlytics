package drip

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jipiboily/forwardlytics/integrations"
)

// Drip integration
type Drip struct {
	api service
}

type service interface {
	request(string, string, []byte) error
}

type dripAPIProduction struct {
	Url string
}

type apiSubscriber struct {
	CustomFields map[string]interface{} `json:"custom_fields"`
	Email        string                 `json:"email"`
	UserId       string                 `json:"user_id"`
}

type apiEvent struct {
	Action     string                 `json:"action"`
	Email      string                 `json:"email"`
	OccurredAt string                 `json:"occurred_at"`
	Properties map[string]interface{} `json:"properties"`
}

type dripAPIError struct {
	Code      string `json:"code"`
	Attribute string `json:"attribute"`
	Message   string `json:"message"`
}

type dripAPIResult struct {
	Errors []dripAPIError `json:"errors"`
}

// Identify forwards and identify call to Drip
func (d Drip) Identify(identification integrations.Identification) (err error) {
	s := apiSubscriber{}
	// Drip needs an email to identify the user
	if identification.UserTraits["email"] == nil {
		logrus.WithField("identification", identification).Error("Drip: Required field email is not present")
		return errors.New("Email is required for doing a drip request")
	} else {
		s.Email = identification.UserTraits["email"].(string)
	}

	s.UserId = string(identification.UserID)

	// Add custom attributes
	s.CustomFields = identification.UserTraits
	s.CustomFields["forwardlyticsReceivedAt"] = identification.ReceivedAt
	s.CustomFields["forwardlyticsTimestamp"] = identification.Timestamp

	payload, err := json.Marshal(map[string][]apiSubscriber{"subscribers": []apiSubscriber{s}})
	err = d.api.request("POST", "subscribers", payload)
	return
}

// Track forwards the event to Drip
func (d Drip) Track(event integrations.Event) (err error) {
	if event.Properties["email"] == nil {
		logrus.WithError(err).WithField("event", event).Error("Drip: Required field email is not present")
		return errors.New("Email is required for doing a drip request")
	}
	e := apiEvent{}
	e.Email = event.Properties["email"].(string)
	event.Properties["forwardlyticsReceivedAt"] = event.ReceivedAt
	e.Action = event.Name
	e.OccurredAt = time.Unix(event.Timestamp, 0).Format("2006-01-02T15:04:05-0700")
	e.Properties = event.Properties
	payload, err := json.Marshal(map[string][]apiEvent{"events": []apiEvent{e}})
	if err != nil {
		logrus.WithField("err", err).Fatal("Error marshalling drip event to json")
	}
	err = d.api.request("POST", "events", payload)
	return
}

// Page forwards the page-events to Drip
// In the drip integration, page-views are just special case events
func (d Drip) Page(page integrations.Page) (err error) {
	if page.Properties["email"] == nil {
		logrus.WithError(err).WithField("page", page).Error("Drip: Required field email is not present")
		return errors.New("Email is required for doing a drip request")
	}
	e := apiEvent{}
	e.Email = page.Properties["email"].(string)
	page.Properties["forwardlyticsReceivedAt"] = page.ReceivedAt
	page.Properties["url"] = page.Url
	page.Properties["pagename"] = page.Name
	e.Action = "Page visited"
	e.OccurredAt = time.Unix(page.Timestamp, 0).Format("2006-01-02T15:04:05-0700")
	e.Properties = page.Properties
	payload, err := json.Marshal(map[string][]apiEvent{"events": []apiEvent{e}})
	if err != nil {
		logrus.WithField("err", err).Fatal("Error marshalling drip page-event to json")
	}
	err = d.api.request("POST", "events", payload)
	if err != nil {
		logrus.WithField("err", err).Fatal("Error from the Drip API...")
	}
	return
}

// Enabled returns wether or not the Drip integration is enabled/configured
func (Drip) Enabled() bool {
	return apiToken() != "" && accountID() != ""
}

func (api dripAPIProduction) request(method string, endpoint string, payload []byte) (err error) {
	apiUrl := api.Url + endpoint
	req, err := http.NewRequest(method, apiUrl, bytes.NewBuffer(payload))
	req.SetBasicAuth(apiToken(), "")
	req.Header.Add("User-Agent", "forwardlytics")
	req.Header.Set("Content-Type", "application/vnd.api+json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logrus.WithError(err).WithField("method", method).WithField("endpoint", endpoint).WithField("payload", string(payload[:])).Error("Error sending request to Drip api")
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.WithError(err).WithField("method", method).WithField("endpoint", endpoint).WithField("payload", string(payload[:])).Error("Error reading body in Drip response")
			return err
		}

		var apiResult dripAPIResult
		json.Unmarshal(body, &apiResult)

		errorDetails := fmt.Sprintf("[%s] %s (on attribute: %s)", apiResult.Errors[0].Code, apiResult.Errors[0].Message, apiResult.Errors[0].Attribute)
		errorMessage := "Drip API returned errors: " + errorDetails

		logrus.WithField("method", method).WithField("endpoint", endpoint).WithField("payload", string(payload[:])).WithFields(
			logrus.Fields{
				"response":    string(body),
				"HTTP-status": resp.StatusCode}).Error(errorMessage)
		return errors.New(errorMessage)
	}
	return
}

func apiUrl() string {
	return "https://api.getdrip.com/v2/" + accountID() + "/"
}

func apiToken() string {
	return os.Getenv("DRIP_API_TOKEN")
}

func accountID() string {
	return os.Getenv("DRIP_ACCOUNT_ID")
}

func init() {
	drip := Drip{}
	drip.api = &dripAPIProduction{Url: apiUrl()}
	integrations.RegisterIntegration("drip", drip)
}
