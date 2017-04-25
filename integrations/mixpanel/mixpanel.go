package mixpanel

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/jipiboily/forwardlytics/integrations"
)

// Mixpanel integration
type Mixpanel struct {
	api service
}

type mixpanelAPIProduction struct {
	Url string
}

type service interface {
	request(string, string, []byte) error
}

type apiSubscriber struct {
	CustomFields map[string]interface{} `json:"$set"`
	UserId       string                 `json:"$distinct_id"`
	Token        string                 `json:"$token"`
	Name         string                 `json:"$name"`
	Email        string                 `json:"$email"`
}

type apiEvent struct {
	Event      string                 `json:"event"`
	Properties map[string]interface{} `json:"properties"`
}

// Identify forwards and identify call to Mixpanel
func (m Mixpanel) Identify(identification integrations.Identification) (err error) {
	s := apiSubscriber{}
	s.UserId = string(identification.UserID)
	s.Token = token()

	if email, ok := identification.UserTraits["email"]; ok {
		s.Email = email.(string)
		delete(identification.UserTraits, "email")
	}

	if name, ok := identification.UserTraits["name"]; ok {
		s.Name = name.(string)
		delete(identification.UserTraits, "name")
	}

	// Add custom attributes
	s.CustomFields = identification.UserTraits
	s.CustomFields["forwardlyticsReceivedAt"] = identification.ReceivedAt
	s.CustomFields["forwardlyticsTimestamp"] = identification.Timestamp

	payload, err := json.Marshal(s)
	err = m.api.request("GET", "engage", payload)
	return
}

// Track forwards the event to Mixpanel
func (m Mixpanel) Track(event integrations.Event) (err error) {
	// Events that are more than 5 years ago will fail in the mixpanel api
	fiveYearsAgo := time.Now().AddDate(-5, 0, 0).Unix()
	if event.Timestamp < fiveYearsAgo {
		logrus.Error("Mixpanel won't accept events that are older than 5 years")
		return errors.New("Mixpanel doesn't support importing events with a timestamp more than 5 years ago")
	}

	e := apiEvent{}
	e.Event = event.Name
	event.Properties["forwardlyticsReceivedAt"] = event.ReceivedAt
	event.Properties["time"] = event.Timestamp
	event.Properties["token"] = token()
	event.Properties["distinct_id"] = event.UserID
	delete(event.Properties, "email")
	e.Properties = event.Properties
	payload, err := json.Marshal(e)
	if err != nil {
		logrus.WithField("err", err).Fatal("Error marshalling Mixpanel event to json")
		return
	}

	// Strange design choice in the mixpanel api: Need to use a
	// different endpoint if the timestamp is more than 5 days
	// ago. Submitting all events older than 5 days to the
	// 'import'-endpoint (see:
	// http://mixpanel.com/help/reference/http "Tracking via
	// HTTP")
	fiveDaysAgo := time.Now().AddDate(0, 0, -5).Unix()
	var endpoint string
	if event.Timestamp > fiveDaysAgo {
		endpoint = "track"
	} else {
		endpoint = "import"
	}

	err = m.api.request("GET", endpoint, payload)
	return
}

// Page forwards the page-events to Mixpanel
func (m Mixpanel) Page(page integrations.Page) (err error) {
	// Mixpanel has no special api for Page events, so just pass on the Page-event to Track
	// Name is the name of the page
	pageEvent := integrations.Event{}
	pageEvent.Name = page.Name
	pageEvent.UserID = page.UserID
	pageEvent.Properties = page.Properties
	pageEvent.Properties["url"] = page.Url
	pageEvent.Properties["event"] = "page"
	pageEvent.Timestamp = page.Timestamp
	pageEvent.ReceivedAt = page.ReceivedAt
	err = m.Track(pageEvent)
	return
}

// Enabled returns wether or not the Mixpanel integration is enabled/configured
func (Mixpanel) Enabled() bool {
	return token() != ""
}

func (api mixpanelAPIProduction) request(method string, endpoint string, payload []byte) (err error) {
	apiUrl := api.Url + endpoint
	req, err := http.NewRequest(method, apiUrl, nil)
	// Mixpanel needs the request to be GET http://<api-url>?data=<base64-encoded payload>
	q := req.URL.Query()
	logrus.Info(base64.StdEncoding.EncodeToString(payload))
	q.Add("data", base64.StdEncoding.EncodeToString(payload))

	// Hack, when using the /import-endpoint api_key must be
	// passed as a query param as well
	if endpoint == "import" {
		q.Add("api_key", token())
	}

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		logrus.WithError(err).WithField("method", method).WithField("endpoint", endpoint).WithField("payload", string(payload[:])).Error("Error sending request to Mixpanel api")
		return
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.WithError(err).WithField("method", method).WithField("endpoint", endpoint).WithField("payload", string(payload[:])).Error("Error reading body in Mixpanel response")
		return err
	}

	// Mixpanel returns a 200OK with a body == 0 when things go wrong
	if resp.StatusCode != http.StatusOK || string(responseBody) == "0" {
		logrus.WithField("method", method).WithField("endpoint", endpoint).WithField("payload", string(payload[:])).WithFields(
			logrus.Fields{
				"response":    string(responseBody),
				"HTTP-status": resp.StatusCode}).Error("Mixpanel api returned errors")
	}

	logrus.WithField("method", method).WithField("endpoint", endpoint).WithField("payload", string(payload[:])).WithFields(
		logrus.Fields{
			"response":    string(responseBody),
			"HTTP-status": resp.StatusCode}).Info("Sent stuff to Mixpanel ", string(req.URL.String()))

	return
}

func token() string {
	return os.Getenv("MIXPANEL_TOKEN")
}

func init() {
	mixpanel := Mixpanel{}
	mixpanel.api = &mixpanelAPIProduction{Url: "http://api.mixpanel.com/"}
	integrations.RegisterIntegration("mixpanel", mixpanel)
}
