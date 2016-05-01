package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jipiboily/forwardlytics/handlers"
	_ "github.com/jipiboily/forwardlytics/integrations/drip"
	_ "github.com/jipiboily/forwardlytics/integrations/intercom"
	_ "github.com/jipiboily/forwardlytics/integrations/keen"
	_ "github.com/jipiboily/forwardlytics/integrations/mixpanel"
)

func main() {
	if os.Getenv("FORWARDLYTICS_API_KEY") == "" {
		log.Fatal("You need to set FORWARDLYTICS_API_KEY")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.Handle("/identify", handlers.AuthMiddleware(http.HandlerFunc(handlers.Identify)))
	http.Handle("/track", handlers.AuthMiddleware(http.HandlerFunc(handlers.Track)))
	log.Println("Forwardlytics started on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
