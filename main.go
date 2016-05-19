package main

import (
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/jipiboily/forwardlytics/handlers"
	_ "github.com/jipiboily/forwardlytics/integrations/drip"
	_ "github.com/jipiboily/forwardlytics/integrations/intercom"
	_ "github.com/jipiboily/forwardlytics/integrations/mixpanel"
)

func main() {
	if os.Getenv("FORWARDLYTICS_API_KEY") == "" {
		logrus.Fatal("You need to set FORWARDLYTICS_API_KEY")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	http.Handle("/identify", handlers.AuthMiddleware(http.HandlerFunc(handlers.Identify)))
	http.Handle("/track", handlers.AuthMiddleware(http.HandlerFunc(handlers.Track)))
	logrus.Info("Forwardlytics started on port", port)
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}
