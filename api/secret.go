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

	s := model.Secret{}
	err := d.Decode(&s)

	if err != nil {
		return nil, errors.BadRequest("")
	}

	author, _ := uuid.FromString(c["sub"].(string))

	tx, _ := model.CreateTransaction()

	defer tx.AutoRollback()

	s, err = s.Create(tx, author)

	if err != nil {
		return nil, errors.BadRequest("")
	}

	for i := range s.Keys {

		s.Keys[i].Secret = s
		k, err := s.Keys[i].Create(tx, s)

		if err != nil {
			return nil, errors.BadRequest("")
		}

		// Prevent circular, stack destorying struct.
		s.Keys[i].Secret = model.Secret{}
		s.Keys[i].ID = k.ID

	}

	tx.Commit()

	return s, nil
}
