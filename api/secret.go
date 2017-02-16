package api

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/roccaforte/server/errors"
	"github.com/roccaforte/server/model"
	"github.com/satori/go.uuid"
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
		r.Put("/:secretID", handler(updateSecret).Serve)
		r.Delete("/:secretID", handler(deleteSecret).Serve)
	})

	return r
}

func allSecrets(w http.ResponseWriter, r *http.Request) (content, error) {
	author := bearerTokenSubject(r)
	secrets, err := model.AllSecrets(author)

	if err != nil {
		return nil, errors.InternalServerError(err.Error())
	}

	return secrets, nil
}

func createSecret(w http.ResponseWriter, r *http.Request) (content, error) {
	s := model.Secret{}
	err := decode(r, &s)

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	author := bearerTokenSubject(r)
	s.Author = model.User{ID: author}

	err = s.Validate()

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	s, err = s.Create()

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	return s, nil
}

func updateSecret(w http.ResponseWriter, r *http.Request) (content, error) {
	secretID, err := uuid.FromString(chi.URLParam(r, "secretID"))

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	s := model.Secret{}
	err = decode(r, &s)

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	old, err := model.GetSecret(secretID)

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	author := bearerTokenSubject(r)

	if author != old.Author.ID {
		return nil, errors.Unauthorized("You do not have permission to edit this secret.")
	}

	s.ID = secretID
	s.Author.ID = author
	s.Update()

	return s, err
}

func deleteSecret(w http.ResponseWriter, r *http.Request) (content, error) {
	secretID, err := uuid.FromString(chi.URLParam(r, "secretID"))

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	s, err := model.GetSecret(secretID)

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	author := bearerTokenSubject(r)

	if author != s.Author.ID {
		return nil, errors.Unauthorized("You do not have permission to delete this secret.")
	}

	err = s.Delete()

	if err != nil {
		return nil, errors.InternalServerError("Unable to delete Secret.")
	}

	return s, nil
}
