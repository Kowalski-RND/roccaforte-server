package api

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/pressly/chi"
	"github.com/roccaforte/server/errors"
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
		r.Get("/", handler(allSecrets).Serve)
		r.Post("/", handler(createSecret).Serve)
	})

	return r
}

func allSecrets(w http.ResponseWriter, r *http.Request) (content, error) {
	defer r.Body.Close()

	c := r.Context().Value(ctxJWT).(jwt.MapClaims)

	author, _ := uuid.FromString(c["sub"].(string))
	secrets, err := model.AllSecrets(author)

	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return secrets, nil
}

func createSecret(w http.ResponseWriter, r *http.Request) (content, error) {
	defer r.Body.Close()

	d := json.NewDecoder(r.Body)

	c := r.Context().Value(ctxJWT).(jwt.MapClaims)

	var s model.Secret
	err := d.Decode(&s)

	if err != nil {
		return nil, errors.BadRequest("")
	}

	author, _ := uuid.FromString(c["sub"].(string))

	s, err = s.Create(author)

	if err != nil {
		return nil, errors.BadRequest("")
	}

	return s, nil
}
