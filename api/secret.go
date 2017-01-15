package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/roccaforte/server/model"
	"github.com/satori/go.uuid"
	"net/http"
)

const (
	ctxJWT contextKey = iota
)

func secretRouter() http.Handler {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(bearerTokenCtx)
		r.Get("/", allSecrets)
		r.Post("/", createSecret)
	})

	return r
}

func allSecrets(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, "Not implemented")
}

func createSecret(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	d := json.NewDecoder(r.Body)

	c := r.Context().Value(ctxJWT).(jwt.MapClaims)

	var s model.Secret
	err := d.Decode(&s)

	if err != nil {
		BadRequest(w, r, "")
		return
	}

	author, err := uuid.FromString(c["sub"].(string))

	if err != nil {
		InternalServerError(w, r, "")
	}

	s.Author = author

	s, err = s.Create()

	if err != nil {
		BadRequest(w, r, "")
		return
	}

	render.JSON(w, r, s)
}
