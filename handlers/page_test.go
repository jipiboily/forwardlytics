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

func TestPageWhenNotPOST(t *testing.T) {
	expectedStatusCode := 404
	expectedBody := "404 page not found"

	r, err := http.NewRequest("GET", "/page", nil)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Page(w, r)

	if w.Code != expectedStatusCode {
		t.Errorf(`Wrong status code. Expecting %v but got %v`, expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

func TestPageWhenInvalidJSON(t *testing.T) {
	expectedStatusCode := 400
	expectedBody := `{"message": "Invalid request."}`

	requestBody := `invalid JSON here`
	r, err := http.NewRequest("POST", "/page", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Page(w, r)

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

func TestPageWhenMissingParameter(t *testing.T) {
	expectedStatusCode := 400
	expectedBody := `{"message": "Missing parameters: name, url, userID, timestamp."}`

	requestBody := `{}`
	r, err := http.NewRequest("POST", "/page", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	Page(w, r)

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

func TestPageWhenOneIntegrationFails(t *testing.T) {
	expectedStatusCode := 500
	expectedBody := `{"message": "Fatal error during page with an integration (test-only-integration-failing): some random error"}`

	requestBody := `{
		"name":"something.created",
		"userID":"123",
                "url":"http://www.example.com",
		"properties": { "someCounter": 97 },
		"timestamp": 12345678
	}`
	r, err := http.NewRequest("POST", "/page", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	failingIntegration := FailingIntegrationPage{}
	integrations.RegisterIntegration("test-only-integration-failing", failingIntegration)
	defer integrations.RemoveIntegration("test-only-integration-failing")

	workingIntegration := FakeIntegration{}
	integrations.RegisterIntegration("test-only-integration-working", workingIntegration)
	defer integrations.RemoveIntegration("test-only-integration-working")

	Page(w, r)

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

func TestPageWhenValid(t *testing.T) {
	expectedStatusCode := 200
	expectedBody := `{"message": "Forwarding page to integrations."}`

	requestBody := `{
		"name":"something.created",
		"userID":"123",
                "url":"http://www.example.com",
		"properties": { "someCounter": 97 },
		"timestamp": 12345678
	}`
	r, err := http.NewRequest("POST", "/page", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()

	integration := &CalledIntegrationPage{t: t}
	integrations.RegisterIntegration("test-only-integration-called", integration)
	defer integrations.RemoveIntegration("test-only-integration-called")

	Page(w, r)

	if !integration.PageCalled {
		t.Error("Page was not called on the integration")
	}

	if w.Code != expectedStatusCode {
		t.Errorf("Wrong status code. Expecting %v but got %v", expectedStatusCode, w.Code)
	}

	if !strings.Contains(w.Body.String(), expectedBody) {
		t.Errorf(`Wrong response. Expecting "%s" but got "%s"`, expectedBody, w.Body.String())
	}
}

// FailingIntegrationPage is an integration that fails when called
type FailingIntegrationPage struct {
	FakeIntegration
}

// Page is failing in this case
func (fi FailingIntegrationPage) Page(page integrations.Page) error {
	return errors.New("some random error")
}

// Enabled returns true because this failing integraiton is enabled
func (FailingIntegrationPage) Enabled() bool {
	return true
}

type CalledIntegrationPage struct {
	FakeIntegration
	PageCalled bool
	t          *testing.T
}

func (i *CalledIntegrationPage) Page(page integrations.Page) error {
	i.PageCalled = true

	expectedReceivedAtCloseTo := time.Now().Unix() - 5
	if page.ReceivedAt < expectedReceivedAtCloseTo {
		i.t.Errorf("ReceivedAt looks wrong. Expecting something close to %v but got %v", expectedReceivedAtCloseTo, page.ReceivedAt)
	}

	return nil
}

func (i CalledIntegrationPage) Enabled() bool {
	return true
}
