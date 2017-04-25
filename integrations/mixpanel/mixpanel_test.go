package mixpanel

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/jipiboily/forwardlytics/integrations"
)

func TestToken(t *testing.T) {
	m := Mixpanel{}
	if m.Enabled() {
		t.Error("Mixpanel shouldn't be enabled without an api-token from ENV")
	}

	os.Setenv("MIXPANEL_TOKEN", "321")
	if !m.Enabled() {
		t.Error("Mixpanel should be enabled when an api-token is set in ENV")
	}

	if token() != "321" {
		t.Errorf("Error in api-token, expected 321, got: %s", token())
	}
}

func TestIdentify(t *testing.T) {
	os.Setenv("MIXPANEL_TOKEN", "321")
	m := Mixpanel{}
	api := APIMock{Url: "http://www.example.com"}
	m.api = &api
	identification := integrations.Identification{
		UserID: "123",
		UserTraits: map[string]interface{}{
			"email": "john@example.com",
			"name":  "John Candy",
		},
		Timestamp:  1234567,
		ReceivedAt: 8765432,
	}
	err := m.Identify(identification)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "GET" {
		t.Errorf("Expected method to be GET, was: %v", api.Method)
	}

	if api.Endpoint != "engage" {
		t.Errorf("Expected endpoint to be engage, was: %v", api.Endpoint)
	}

	expectedPayload := `{"$set":{"forwardlyticsReceivedAt":8765432,"forwardlyticsTimestamp":1234567},"$distinct_id":"123","$token":"321","$name":"John Candy","$email":"john@example.com"}`
	if string(api.Payload) != expectedPayload {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %s", api.Payload)
	}
}

func TestTrackRecentEvent(t *testing.T) {
	os.Setenv("MIXPANEL_TOKEN", "321")
	m := Mixpanel{}
	api := APIMock{Url: "http://www.example.com"}
	m.api = &api

	timestampForTest := time.Now().Unix()

	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  timestampForTest,
		ReceivedAt: 65,
	}

	err := m.Track(event)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "GET" {
		t.Errorf("Expected method to be GET, was: %v", api.Method)
	}

	if api.Endpoint != "track" {
		t.Errorf("Expected endpoint to be track, was: %v", api.Endpoint)
	}

	expectedPayload := "{\"event\":\"account.created\",\"properties\":{\"distinct_id\":\"123\",\"forwardlyticsReceivedAt\":65,\"time\":" + strconv.Itoa(int(timestampForTest)) + ",\"token\":\"321\"}}"
	if string(api.Payload) != expectedPayload {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %s", api.Payload)
	}
}

func TestTrackEventOlderThanFiveDays(t *testing.T) {
	os.Setenv("MIXPANEL_TOKEN", "321")
	m := Mixpanel{}
	api := APIMock{Url: "http://www.example.com"}
	m.api = &api

	timestampForTest := time.Now().AddDate(0, 0, -6).Unix()

	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  timestampForTest,
		ReceivedAt: 65,
	}

	err := m.Track(event)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "GET" {
		t.Errorf("Expected method to be GET, was: %v", api.Method)
	}

	if api.Endpoint != "import" {
		t.Errorf("Expected endpoint to be import, was: %v", api.Endpoint)
	}

	expectedPayload := "{\"event\":\"account.created\",\"properties\":{\"distinct_id\":\"123\",\"forwardlyticsReceivedAt\":65,\"time\":" + strconv.Itoa(int(timestampForTest)) + ",\"token\":\"321\"}}"
	if string(api.Payload) != expectedPayload {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %s", api.Payload)
	}
}

func TestTrackEventOlderThanFiveYears(t *testing.T) {
	os.Setenv("MIXPANEL_TOKEN", "321")
	m := Mixpanel{}
	api := APIMock{Url: "http://www.example.com"}
	m.api = &api

	timestampForTest := time.Now().AddDate(-6, 0, 0).Unix()

	event := integrations.Event{
		Name:   "account.created",
		UserID: "123",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  timestampForTest,
		ReceivedAt: 65,
	}

	err := m.Track(event)
	if err == nil {
		t.Error("Events with a timestamp more than 5 years ago shouldn't be accepted")
	}
}

// Test page
func TestPage(t *testing.T) {
	os.Setenv("MIXPANEL_TOKEN", "321")
	m := Mixpanel{}
	api := APIMock{Url: "http://www.example.com"}
	m.api = &api

	timestampForTest := time.Now().Unix()

	page := integrations.Page{
		Name:   "Homepage",
		UserID: "123",
		Url:    "http://www.example.com",
		Properties: map[string]interface{}{
			"email": "john@example.com",
		},
		Timestamp:  timestampForTest,
		ReceivedAt: 65,
	}

	err := m.Page(page)
	if err != nil {
		t.Fatal(err)
	}

	if api.Method != "GET" {
		t.Errorf("Expected method to be GET, was: %v", api.Method)
	}

	if api.Endpoint != "track" {
		t.Errorf("Expected endpoint to be track, was: %v", api.Endpoint)
	}

	expectedPayload := "{\"event\":\"Homepage\",\"properties\":{\"distinct_id\":\"123\",\"event\":\"page\",\"forwardlyticsReceivedAt\":65,\"time\":" + strconv.Itoa(int(timestampForTest)) + ",\"token\":\"321\",\"url\":\"http://www.example.com\"}}"
	if string(api.Payload) != expectedPayload {
		t.Errorf("Expected payload: "+string(expectedPayload)+" got: %s", api.Payload)
	}
}

func TestRequest(t *testing.T) {

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
