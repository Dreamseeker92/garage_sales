package web

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// Represents response to the client when an error occurs.
type ErrorResponse struct {
	Error string `json:"error"`
	Fields []FieldError
}

// Provides web information to the request error
type Error struct {
  Err error
  Status int
  Fields []FieldError
}

func NewRequestError(err error, status int) *Error {
	return &Error{Err: err, Status: status}
}

func (e *Error) Error() string {
	return e.Err.Error()
}
