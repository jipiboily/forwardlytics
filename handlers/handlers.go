package handlers

import (
	"fmt"
	"net/http"
)

func writeResponse(w http.ResponseWriter, body string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	body = fmt.Sprintf(`{"message": "%s"}`, body)
	w.Write([]byte(body))

}
