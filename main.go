package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jipiboily/forwardlytics/integrations"

	_ "github.com/jipiboily/forwardlytics/integrations/drip"
	_ "github.com/jipiboily/forwardlytics/integrations/intercom"
	_ "github.com/jipiboily/forwardlytics/integrations/keen"
	_ "github.com/jipiboily/forwardlytics/integrations/mixpanel"
)

func main() {
	if os.Getenv("FORWARDLYTICS_API_KEY") == "" {
		log.Fatal("You need to set FORWARDLYTICS_API_KEY")
	}

	http.HandleFunc("/identify", identifyHandler)
	log.Println("Forwardlytics started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func identifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	apiKey := r.Header.Get("FORWARDLYTICS_API_KEY")
	if apiKey != os.Getenv("FORWARDLYTICS_API_KEY") {
		log.Printf("Wrong API key. We had '%s' but it should be '%s'\n", apiKey, os.Getenv("FORWARDLYTICS_API_KEY"))

		errorMsg := "Invalid API KEY. The FORWARDLYTICS_API_KEY header must be specified, with the proper API key."
		writeResponse(w, errorMsg, http.StatusUnauthorized)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var user integrations.User
	err := decoder.Decode(&user)
	if err != nil {
		log.Println("Bad request:", r.Body)
		writeResponse(w, "Invalid request.", http.StatusBadRequest)
		return
	}

	for _, integrationName := range integrations.IntegrationList() {
		integration := integrations.GetIntegration(integrationName)
		if integration.Enabled() {
			log.Println("Forwarding idenitify to", integrationName)
			integration.Identify(user)
		}

	}

	writeResponse(w, "Forwarding identify to integrations.", http.StatusOK)
}

func writeResponse(w http.ResponseWriter, body string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	body = fmt.Sprintf(`{"message": "%s"}`, body)
	w.Write([]byte(body))

}
