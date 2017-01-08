package api

import (
	"github.com/pressly/chi/render"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, msg string) {
	render.Status(r, http.StatusInternalServerError)
	if msg != "" {
		render.JSON(w, r, msg)
	} else {
		render.JSON(w, r, http.StatusText(http.StatusInternalServerError))
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request, msg string) {
	render.Status(r, http.StatusBadRequest)
	if msg != "" {
		render.JSON(w, r, msg)
	} else {
		render.JSON(w, r, http.StatusText(http.StatusBadRequest))
	}
}

func NotFound(w http.ResponseWriter, r *http.Request, msg string) {
	render.Status(r, http.StatusNotFound)
	if msg != "" {
		render.JSON(w, r, msg)
	} else {
		render.JSON(w, r, http.StatusText(http.StatusNotFound))
	}
}
