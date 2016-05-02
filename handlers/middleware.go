package handlers

import (
	"net/http"
	"os"
)

// AuthMiddleware is making sure the call is properly authenticated before
// sending the request to the handlers
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("FORWARDLYTICS_API_KEY")
		if apiKey != os.Getenv("FORWARDLYTICS_API_KEY") {
			errorMsg := "Invalid API KEY. The FORWARDLYTICS_API_KEY header must be specified, with the proper API key."
			writeResponse(w, errorMsg, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
