package errors

import (
	"github.com/go-errors/errors"
	"net/http"
)

type Error struct {
	err    error
	status int
}

func (e Error) Error() string {
	return e.Error()
}

func (e Error) Status() int {
	return e.status
}

func New(msg string, status int) Error {
	return Error{
		err:    errors.Errorf(msg),
		status: status,
	}
}

func BadRequest(msg string) Error {
	return New(msg, http.StatusBadRequest)
}

func InternalServerError(msg string) Error {
	return New(msg, http.StatusInternalServerError)
}

func NotFound(msg string) Error {
	return New(msg, http.StatusNotFound)
}

func Unauthorized(msg string) Error {
	return New(msg, http.StatusUnauthorized)
}
