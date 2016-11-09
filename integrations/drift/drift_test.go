package drift

import (
	"os"
	"testing"

	"github.com/jipiboily/forwardlytics/integrations"
)

func TestNotEnabledWhenMissingCredentials(t *testing.T) {
	os.Setenv("DRIFT_ORG_ID", "")
	drift := Drift{}
	if drift.Enabled() {
		t.Error("Should not be enabled when missing accoundId and apiToken")
	}
}

func TestEnabledWhenCredentialsPresent(t *testing.T) {
	os.Setenv("DRIFT_ORG_ID", "123")
	drift := Drift{}
	if !drift.Enabled() {
		t.Error("Should be enabled when org-id present")
	}
}

func TestIdentify(t *testing.T) {
	os.Setenv("DRIFT_ORG_ID", "123")
	drift := Drift{}
	api := APIMock{baseUrl: "http://www.example.com"}
	drift.api = &api
	identification := integrations.Identification{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 8765432,
	}
	err := drift.Identify(identification)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "POST" {
		t.Errorf("Expected method to be POST, was: %v", api.Method)
	}

	if api.Endpoint != "identify" {
		t.Errorf("Expected endpoint to be identify, was: %v", api.Endpoint)
	}

	expectedPayload := `{"attributes":{"email":"john@example.com","forwardlyticsReceivedAt":8765432,"forwardlyticsTimestamp":1234567},"createdAt":1234567,"userId":"123","orgId":"123"}`
	if string(api.Payload) != expectedPayload {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %s", api.Payload)
	}
}

func TestTrack(t *testing.T) {
	os.Setenv("DRIFT_ORG_ID", "321")
	drift := Drift{}
	api := APIMock{baseUrl: "http://www.example.com/"}
	drift.api = &api
	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  1234567,
		ReceivedAt: 65,
	}

	err := drift.Track(event)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "POST" {
		t.Errorf("Expected method to be POST, was: %v", api.Method)
	}

	if api.Endpoint != "track" {
		t.Errorf("Expected endpoint to be events, was: %v", api.Endpoint)
	}
	expectedPayload := `{"orgId":"321","userId":"123","event":"account.created","createdAt":1234567,"attributes":{"email":"john@example.com","forwardlyticsReceivedAt":65}}`

	if string(api.Payload) != string(expectedPayload) {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %v", string(api.Payload))
	}
}

func TestPage(t *testing.T) {
	os.Setenv("DRIFT_ORG_ID", "321")
	drift := Drift{}
	api := APIMock{baseUrl: "http://www.example.com/"}
	drift.api = &api
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

	err := drift.Page(page)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "POST" {
		t.Errorf("Expected method to be POST, was: %v", api.Method)
	}

	if api.Endpoint != "track" {
		t.Errorf("Expected endpoint to be track, was: %v", api.Endpoint)
	}
	expectedPayload := `{"orgId":"321","userId":"123","event":"page","url":"http://www.example.com","createdAt":1234567,"attributes":{"email":"john@example.com","forwardlyticsReceivedAt":65,"name":"Homepage"}}`

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

type APIMock struct {
	baseUrl  string
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
