package api

import (
	"github.com/pressly/chi/render"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, msg string) {
	response(w, r, msg, http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	response(w, r, msg, http.StatusBadRequest)
}

func NotFound(w http.ResponseWriter, r *http.Request, msg string) {
	response(w, r, msg, http.StatusNotFound)
}

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
