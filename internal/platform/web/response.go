package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Marshals a value to JSON and sends it to the client.
func Respond(writer http.ResponseWriter, value interface{}, statusCode int) error {
	data, err := json.Marshal(value)
	if err != nil {
		return errors.Wrap(err, "Marshalling value to JSON")
	}

	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(statusCode)
	if _, err := writer.Write(data); err != nil {
		return errors.Wrap(err, "Writing to the client")
	}

	return nil
}

// Knows how to handle errors going out to the client.
func RespondError(writer http.ResponseWriter, err error) error {

	if webErr, ok := err.(*Error); ok {

		body := ErrorResponse{Error: err.Error()}
		return Respond(writer, body, webErr.Status)
	}

	return Respond(writer, ErrorResponse{Error: http.StatusText(http.StatusInternalServerError)}, http.StatusInternalServerError)
}
