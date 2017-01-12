package api

import (
	"github.com/pressly/chi/render"
	"net/http"
)

// InternalServerError returns a JSON message with a 500 status code.
// If msg is empty string then the message sent is the status text
// from the net/http library corresponding to the status code.
func InternalServerError(w http.ResponseWriter, r *http.Request, msg string) {
	response(w, r, msg, http.StatusInternalServerError)
}

// BadRequest returns a JSON message with a 400 status code.
// If msg is empty string then the message sent is the status text
// from the net/http library corresponding to the status code.
func BadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	response(w, r, msg, http.StatusBadRequest)
}

// NotFound returns a JSON message with a 404 status code.
// If msg is empty string then the message sent is the status text
// from the net/http library corresponding to the status code.
func NotFound(w http.ResponseWriter, r *http.Request, msg string) {
	response(w, r, msg, http.StatusNotFound)
}

// Unauthorized returns a JSON message with a 401 status code.
// If msg is empty string then the message sent is the status text
// from the net/http library corresponding to the status code.
func Unauthorized(w http.ResponseWriter, r *http.Request, msg string) {
	response(w, r, msg, http.StatusUnauthorized)
}

func response(w http.ResponseWriter, r *http.Request, msg string, status int) {
	render.Status(r, status)
	if msg != "" {
		render.JSON(w, r, msg)
	} else {
		render.JSON(w, r, http.StatusText(status))
	}
}
