package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/codeship/go-retro"
)

func resourceNotReady(resourceError error) error {
	if os.Getenv("NUM_RETRIES_ON_ERROR") == "" {
		return resourceError
	}
	numRetries, err := strconv.Atoi(os.Getenv("NUM_RETRIES_ON_ERROR"))
	if err != nil {
		logrus.WithField("err", err).Error("env variable NUM_RETRIES_ON_ERROR should be an integer")
		return err
	}
	logrus.WithField("error", resourceError).Error("Error sending request")
	return retro.NewBackoffRetryableError(resourceError, numRetries)
}

func writeResponse(w http.ResponseWriter, body string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	body = fmt.Sprintf(`{"message": "%s"}`, body)
	w.Write([]byte(body))

}
