package handlers

import (
	"crypto/subtle"
	"net/http"
	"os"
)

// AuthMiddleware is making sure the call is properly authenticated before
// sending the request to the handlers
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Forwardlytics-Api-Key")

		if subtle.ConstantTimeCompare([]byte(os.Getenv("FORWARDLYTICS_API_KEY")), []byte(apiKey)) != 1 {
			errorMsg := "Invalid API KEY. The Forwardlytics-Api-Key header must be specified, with the proper API key."
			writeResponse(w, errorMsg, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
