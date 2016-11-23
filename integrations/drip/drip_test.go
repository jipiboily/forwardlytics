package drip

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/jipiboily/forwardlytics/integrations"
)

func TestIdentifyErrorWhenNoEmail(t *testing.T) {
	drip := Drip{}
	identification := integrations.Identification{
		UserID:     "123",
		UserTraits: map[string]interface{}{},
	}
	err := drip.Identify(identification)
	if err == nil {
		t.Error("Expected error when no email given")
	}
}

func TestNotEnabledWhenMissingCredentials(t *testing.T) {
	os.Setenv("DRIP_API_TOKEN", "")
	os.Setenv("DRIP_ACCOUNT_ID", "")
	drip := Drip{}
	if drip.Enabled() {
		t.Error("Should not be enabled when missing accoundId and apiToken")
	}
	os.Setenv("DRIP_API_TOKEN", "123")
	if drip.Enabled() {
		t.Error("Should not be enabled when missing accoundId")
	}
	os.Setenv("DRIP_API_TOKEN", "")
	os.Setenv("DRIP_ACCOUNT_ID", "123")
	if drip.Enabled() {
		t.Error("Should not be enabled when missing apiToken")
	}
}

func TestEnabledWhenCredentialsPresent(t *testing.T) {
	os.Setenv("DRIP_API_TOKEN", "123")
	os.Setenv("DRIP_ACCOUNT_ID", "123")
	drip := Drip{}
	if !drip.Enabled() {
		t.Error("Drip should be enabled when credentials are set")
	}
}

func TestIdentify(t *testing.T) {
	os.Setenv("DRIP_API_TOKEN", "123")
	os.Setenv("DRIP_ACCOUNT_ID", "321")
	drip := Drip{}
	api := APIMock{Url: "http://www.example.com"}
	drip.api = &api
	identification := integrations.Identification{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 8765432,
	}
	err := drip.Identify(identification)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "POST" {
		t.Errorf("Expected method to be POST, was: %v", api.Method)
	}

	if api.Endpoint != "subscribers" {
		t.Errorf("Expected endpoint to be subscribers, was: %v", api.Endpoint)
	}

	expectedPayload := `{"subscribers":[{"custom_fields":{"email":"john@example.com","forwardlyticsReceivedAt":8765432,"forwardlyticsTimestamp":1234567},"email":"john@example.com","user_id":"123"}]}`
	if string(api.Payload) != expectedPayload {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %s", api.Payload)
	}
}

func TestTrackErrorWhenNoEmailPresent(t *testing.T) {
	os.Setenv("DRIP_API_TOKEN", "123")
	os.Setenv("DRIP_ACCOUNT_ID", "321")
	api := &APIMock{Url: "http://www.example.com"}
	drip := Drip{}
	drip.api = api
	event := integrations.Event{
		Name:       "account.created",
		UserID:     "123",
		Properties: map[string]interface{}{},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}
	err := drip.Track(event)
	if err == nil {
		t.Error("Expected error when no email given")
	}
}

func TestTrack(t *testing.T) {
	os.Setenv("DRIP_API_TOKEN", "123")
	os.Setenv("DRIP_ACCOUNT_ID", "321")
	drip := Drip{}
	api := APIMock{Url: "http://www.example.com"}
	drip.api = &api
	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}

	err := drip.Track(event)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "POST" {
		t.Errorf("Expected method to be POST, was: %v", api.Method)
	}

	if api.Endpoint != "events" {
		t.Errorf("Expected endpoint to be events, was: %v", api.Endpoint)
	}

	expectedData := &MockEvents{
		Events: []MockEvent{
			{
				Action:     "account.created",
				Email:      "john@example.com",
				OccurredAt: time.Unix(1234567, 0).Format("2006-01-02T15:04:05-0700"),
				Properties: MockEventProperties{
					Email: "john@example.com",
					ForwardlyticsReceivedAt: 65,
				},
			}},
	}

	expectedPayload, _ := json.Marshal(expectedData)
	if string(api.Payload) != string(expectedPayload) {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %v", string(api.Payload))
	}
}

func TestPageErrorWhenNoEmailPresent(t *testing.T) {
	os.Setenv("DRIP_API_TOKEN", "123")
	os.Setenv("DRIP_ACCOUNT_ID", "321")
	api := &APIMock{Url: "http://www.example.com"}
	drip := Drip{}
	drip.api = api
	page := integrations.Page{
		Name:       "Homepage",
		UserID:     "123",
		Url:        "http://www.example.com",
		Properties: map[string]interface{}{},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}
	err := drip.Page(page)
	if err == nil {
		t.Error("Expected error when no email given")
	}
}

func TestPage(t *testing.T) {
	os.Setenv("DRIP_API_TOKEN", "123")
	os.Setenv("DRIP_ACCOUNT_ID", "321")
	drip := Drip{}
	api := APIMock{Url: "http://www.example.com"}
	drip.api = &api
	page := integrations.Page{
		Name:   "Homepage",
		UserID: "123",
		Url:    "http://www.example.com",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}

	err := drip.Page(page)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "POST" {
		t.Errorf("Expected method to be POST, was: %v", api.Method)
	}

	if api.Endpoint != "events" {
		t.Errorf("Expected endpoint to be events, was: %v", api.Endpoint)
	}

	expectedData := &MockPageEvents{
		Events: []MockPageEvent{
			{
				Action:     "Page visited",
				Email:      "john@example.com",
				OccurredAt: time.Unix(1234567, 0).Format("2006-01-02T15:04:05-0700"),
				Properties: MockPageEventProperties{
					Email: "john@example.com",
					ForwardlyticsReceivedAt: 65,
					PageName:                "Homepage",
					Url:                     "http://www.example.com",
				},
			}},
	}

	expectedPayload, _ := json.Marshal(expectedData)
	if string(api.Payload) != string(expectedPayload) {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %v", string(api.Payload))
	}
}

type MockEvents struct {
	Events []MockEvent `json:"events"`
}

type MockEvent struct {
	Action     string              `json:"action"`
	Email      string              `json:"email"`
	OccurredAt string              `json:"occurred_at"`
	Properties MockEventProperties `json:"properties"`
}

type MockEventProperties struct {
	Email                   string `json:"email"`
	ForwardlyticsReceivedAt int64  `json:"forwardlyticsReceivedAt"`
}

type MockPageEvents struct {
	Events []MockPageEvent `json:"events"`
}

type MockPageEvent struct {
	Action     string                  `json:"action"`
	Email      string                  `json:"email"`
	OccurredAt string                  `json:"occurred_at"`
	Properties MockPageEventProperties `json:"properties"`
}

type MockPageEventProperties struct {
	Email                   string `json:"email"`
	ForwardlyticsReceivedAt int64  `json:"forwardlyticsReceivedAt"`
	PageName                string `json:"pagename"`
	Url                     string `json:"url"`
}

type APIMock struct {
	Url      string
	Method   string
	Endpoint string
	Payload  []byte
}

func (api *APIMock) request(method string, endpoint string, payload []byte) error {
	api.Method = method
	api.Endpoint = endpoint
	api.Payload = payload
	return nil
}
