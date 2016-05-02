package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/jipiboily/forwardlytics/integrations"
)

func TestTrackWhenNotPOST(t *testing.T) {
	expectedStatusCode := 404
	expectedBody := "404 page not found"

	r, err := http.NewRequest("GET", "/track", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Track(w, r)

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

func TestTrackWhenInvalidJSON(t *testing.T) {
	expectedStatusCode := 400
	expectedBody := `{"message": "Invalid request."}`

	requestBody := `invalid JSON here`
	r, err := http.NewRequest("POST", "/track", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Track(w, r)

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

func TestTrackWhenMissingParameter(t *testing.T) {
	expectedStatusCode := 400
	expectedBody := `{"message": "Missing parameters: name, userID, timestamp."}`

	requestBody := `{}`
	r, err := http.NewRequest("POST", "/track", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Track(w, r)

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

func TestTrackWhenOneIntegrationFails(t *testing.T) {
	expectedStatusCode := 500
	expectedBody := `{"message": "Fatal error during event with an integration (test-only-integration-failing): some random error"}`

	requestBody := `{
		"name":"something.created",
		"userID":"123",
		"properties": { "someCounter": 97 },
		"timestamp": 12345678
	}`
	r, err := http.NewRequest("POST", "/track", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	failingIntegration := FailingIntegration{}
	integrations.RegisterIntegration("test-only-integration-failing", failingIntegration)
	defer integrations.RemoveIntegration("test-only-integration-failing")

	workingIntegration := FakeIntegration{}
	integrations.RegisterIntegration("test-only-integration-working", workingIntegration)
	defer integrations.RemoveIntegration("test-only-integration-working")

	Track(w, r)

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

func TestTrackWhenValid(t *testing.T) {
	expectedStatusCode := 200
	expectedBody := `{"message": "Forwarding event to integrations."}`

	requestBody := `{
		"name":"something.created",
		"userID":"123",
		"properties": { "someCounter": 97 },
		"timestamp": 12345678
	}`
	r, err := http.NewRequest("POST", "/track", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	integration := &CalledIntegration{t: t}
	integrations.RegisterIntegration("test-only-integration-called", integration)
	defer integrations.RemoveIntegration("test-only-integration-called")

	Track(w, r)

	if !integration.Tracked {
		t.Error("Track was not called on the integration")
	}

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

// FailingIntegration is an integration that fails when called
type FailingIntegration struct {
	FakeIntegration
}

// Track is failing in this case
func (fi FailingIntegration) Track(event integrations.Event) error {
	return errors.New("some random error")
}

// Enabled returns true because this failing integraiton is enabled
func (FailingIntegration) Enabled() bool {
	return true
}

type CalledIntegration struct {
	FakeIntegration
	Tracked bool
	t       *testing.T
}

func (i *CalledIntegration) Track(event integrations.Event) error {
	i.Tracked = true

	expectedReceivedAtCloseTo := time.Now().Unix() - 5
	if event.ReceivedAt < expectedReceivedAtCloseTo {
		i.t.Errorf("ReceivedAt looks wrong. Expecting something close to %v but got %v", expectedReceivedAtCloseTo, event.ReceivedAt)
	}

	return nil
}

func (i CalledIntegration) Enabled() bool {
	return true
}
