package web

import (
	"encoding/json"
	"github.com/pkg/errors"
	"net/http"
)

// Decodes request body into dest.
func Decode(request *http.Request, dest interface{}) error {
	if err := json.NewDecoder(request.Body).Decode(dest); err != nil {
		return errors.Wrap(err, "Decoding request body")
	}
	
	return nil
}
