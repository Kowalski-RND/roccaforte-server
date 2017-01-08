package api

import (
	"github.com/pressly/chi/render"
	"net/http"
)

func InternalServerError(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, http.StatusText(http.StatusInternalServerError))
}

func BadRequest(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, http.StatusText(http.StatusBadRequest))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, http.StatusText(http.StatusNotFound))
}
