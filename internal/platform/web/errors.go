package web

// Represents response to the client when an error occurs.
type ErrorResponse struct {
	Error string `json:"error"`
}
