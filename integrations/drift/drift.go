package drift

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jipiboily/forwardlytics/integrations"
)

// Drift integration
type Drift struct {
	api service
}

type service interface {
	request(string, string, []byte) error
}

type driftAPIProduction struct {
	baseUrl string
}

type apiSubscriber struct {
	Attributes map[string]interface{} `json:"attributes"`
	CreatedAt  int64                  `json:"createdAt"`
	UserId     string                 `json:"userId"`
	OrgId      string                 `json:"orgId"`
}

type apiEvent struct {
	OrgId      string                 `json:"orgId"`
	UserId     string                 `json:"userId"`
	Event      string                 `json:"event"`
	CreatedAt  int64                  `json:"createdAt"`
	Attributes map[string]interface{} `json:"attributes"`
}

type apiPage struct {
	OrgId      string                 `json:"orgId"`
	UserId     string                 `json:"userId"`
	Event      string                 `json:"event"`
	Url        string                 `json:"url"`
	CreatedAt  int64                  `json:"createdAt"`
	Attributes map[string]interface{} `json:"attributes"`
}

// Identify forwards and identify call to Drift
func (d Drift) Identify(identification integrations.Identification) (err error) {
	s := apiSubscriber{}
	s.UserId = string(identification.UserID)
	s.CreatedAt = identification.Timestamp
	s.OrgId = orgID()
	// Add custom attributes
	s.Attributes = identification.UserTraits
	s.Attributes["forwardlyticsReceivedAt"] = identification.ReceivedAt
	s.Attributes["forwardlyticsTimestamp"] = identification.Timestamp
	payload, err := json.Marshal(s)

	err = d.api.request("POST", "identify", payload)
	if err != nil {
		logrus.WithError(err).WithField("identify", identification).WithField("payload", payload).Error("Error sending identify to drift")
	}
	return
}

// Track forwards the event to Drift
func (d Drift) Track(event integrations.Event) (err error) {
	e := apiEvent{}
	e.OrgId = orgID()
	e.UserId = event.UserID
	e.Attributes = event.Properties
	event.Properties["forwardlyticsReceivedAt"] = event.ReceivedAt
	e.Event = event.Name
	e.CreatedAt = event.Timestamp
	payload, err := json.Marshal(e)
	if err != nil {
		logrus.WithError(err).WithField("event", event).WithField("payload", payload).Error("Error marshalling drift event to json")
	}
	err = d.api.request("POST", "track", payload)
	if err != nil {
		logrus.WithError(err).WithField("event", event).WithField("payload", payload).Error("Error sending event to drift")
	}
	return
}

// Page forwards the page-events to Drift
func (d Drift) Page(page integrations.Page) (err error) {
	p := apiPage{}
	p.OrgId = orgID()
	p.UserId = page.UserID
	p.Url = page.Url
	page.Properties["forwardlyticsReceivedAt"] = page.ReceivedAt
	page.Properties["name"] = page.Name
	p.Attributes = page.Properties
	p.Event = "page"
	p.CreatedAt = page.Timestamp
	payload, err := json.Marshal(p)
	if err != nil {
		logrus.WithError(err).WithField("page", page).WithField("payload", payload).Error("Error marshalling drift page-event to json")
	}
	err = d.api.request("POST", "track", payload)
	if err != nil {
		logrus.WithError(err).WithField("page", page).WithField("payload", payload).Error("Error sending page-event to drift")
	}
	return
}

// Enabled returns wether or not the Drift integration is enabled/configured
func (Drift) Enabled() bool {
	return orgID() != ""
}

func (api driftAPIProduction) request(method string, endpoint string, payload []byte) (err error) {
	apiUrl := api.baseUrl + endpoint
	req, err := http.NewRequest(method, apiUrl, bytes.NewBuffer(payload))
	req.Header.Add("User-Agent", "forwardlytics")
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logrus.WithError(err).WithFields(
			logrus.Fields{
				"method":   method,
				"apiUrl":   apiUrl,
				"endpoint": endpoint,
				"payload":  payload}).Error("Error sending request to Drift api")
		return
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.WithError(err).WithFields(
				logrus.Fields{
					"method":     method,
					"apiUrl":     apiUrl,
					"endpoint":   endpoint,
					"payload":    payload,
					"httpstatus": resp.StatusCode}).Error("Error reading Drift response")
			return err
		}
		logrus.WithFields(
			logrus.Fields{
				"response":    string(body),
				"HTTP-status": resp.StatusCode,
				"method":      method,
				"endpoint":    endpoint,
				"payload":     payload}).Error("Drift api returned errors")

	}
	return
}

func orgID() string {
	return os.Getenv("DRIFT_ORG_ID")
}

func init() {
	drift := Drift{}
	drift.api = &driftAPIProduction{baseUrl: "https://event.api.drift.com/"}
	integrations.RegisterIntegration("drift", drift)
}
