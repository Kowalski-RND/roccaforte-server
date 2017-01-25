package api

import (
	"encoding/json"
	"github.com/pressly/chi"
	"github.com/roccaforte/server/errors"
	"github.com/roccaforte/server/model"
	"net/http"
)

const (
	keyUser contextKey = iota
)

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", handler(allUsers).Serve)
	r.Post("/", handler(createUser).Serve)

	r.Route("/:username", func(r chi.Router) {
		r.Use(bearerTokenCtx)
		r.Get("/", handler(getUser).Serve)
	})

	return r
}

func allUsers(w http.ResponseWriter, r *http.Request) (content, error) {
	defer r.Body.Close()

	users, err := model.AllUsers()

	if err != nil {
		return nil, errors.InternalServerError("")
	}

	return users, nil
}

func createUser(w http.ResponseWriter, r *http.Request) (content, error) {
	defer r.Body.Close()

	d := json.NewDecoder(r.Body)

	var u model.User
	err := d.Decode(&u)

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	err = u.Create()

	if err != nil {
		return nil, errors.BadRequest(err.Error())
	}

	return u, nil
}

func getUser(w http.ResponseWriter, r *http.Request) (content, error) {
	defer r.Body.Close()

	username := chi.URLParam(r, "username")
	u, err := model.UserByUsername(username)

	if err != nil {
		return nil, errors.InternalServerError("")
	} else if (model.User{}) == u {
		return nil, errors.NotFound("No user found for given username.")
	}

	return u, nil
}
