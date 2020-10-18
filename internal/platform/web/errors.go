package web

// Represents response to the client when an error occurs.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Provides web information to the request error
type Error struct {
  Err error
  Status int
}

func NewRequestError(err error, status int) *Error {
	return &Error{Err: err, Status: status}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
