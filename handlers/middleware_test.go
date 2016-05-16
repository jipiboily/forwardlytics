package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAuthMiddlewareWhenAPIKeyIsValid(t *testing.T) {
	os.Setenv("FORWARDLYTICS_API_KEY", "DNUAS67AASNDj")
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Forwardlytics-Api-Key", `DNUAS67AASNDj`)
	rr := httptest.NewRecorder()

	AuthMiddleware(FakeHandler{}).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"message": "success"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAuthMiddlewareWhenAPIKeyIsInvalid(t *testing.T) {
	os.Setenv("FORWARDLYTICS_API_KEY", "DNUAS67AASNDj")
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Forwardlytics-Api-Key", `...`)
	rr := httptest.NewRecorder()

	AuthMiddleware(FakeHandler{}).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	expected := `{"message": "Invalid API KEY. The Forwardlytics-Api-Key header must be specified, with the proper API key."}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

type FakeHandler struct{}

func (FakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, "success", 200)
}
