package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/codeship/go-retro"
)

var ErrNotReady = retro.NewBackoffRetryableError(errors.New("error: resource not ready"), 10)

func writeResponse(w http.ResponseWriter, body string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	body = fmt.Sprintf(`{"message": "%s"}`, body)
	w.Write([]byte(body))

}
