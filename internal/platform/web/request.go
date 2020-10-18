package web

import (
	"encoding/json"
	"net/http"
)

// Decodes request body into dest.
func Decode(request *http.Request, dest interface{}) error {
	if err := json.NewDecoder(request.Body).Decode(dest); err != nil {
		return NewRequestError(err, http.StatusBadRequest)
	}
	
	return nil
}
