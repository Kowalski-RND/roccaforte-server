package errors

import (
	"net/http"
)

// Error captures a message along with a HTTP Status code to
// be used by the API handler.
type Error struct {
	msg    string
	status int
}

// Error is present to implement the Error interface.
func (e Error) Error() string {
	return e.msg
}

// Status returns the HTTP Status code associated with this Error.
func (e Error) Status() int {
	return e.status
}

// New creates a new Error for given message and status.
func New(msg string, status int) Error {
	return Error{
		msg:    msg,
		status: status,
	}
}

// BadRequest returns an error with a given message and http.StatusBadRequest
// as the status code.
func BadRequest(msg string) Error {
	return New(msg, http.StatusBadRequest)
}

// InternalServerError returns an error with a given message and http.StatusInternalServerError
// as the status code.
func InternalServerError(msg string) Error {
	return New(msg, http.StatusInternalServerError)
}

// NotFound returns an error with a given message and http.StatusNotFound
// as the status code.
func NotFound(msg string) Error {
	return New(msg, http.StatusNotFound)
}

// Unauthorized returns an error with a given message and http.StatusUnauthorized
// as the status code.
func Unauthorized(msg string) Error {
	return New(msg, http.StatusUnauthorized)
}
