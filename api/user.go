package api

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/roccaforte/server/model"
	"net/http"
)

const (
	keyUser contextKey = iota
)

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", allUsers)
	r.Post("/", createUser)

	r.Route("/:username", func(r chi.Router) {
		r.Use(bearerTokenCtx)
		r.Use(userCtx)
		r.Get("/", getUser)
	})

	return r
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	users, err := model.AllUsers()

	if err != nil {
		InternalServerError(w, r, "")
		return
	}

	render.JSON(w, r, users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	d := json.NewDecoder(r.Body)

	var u model.User
	err := d.Decode(&u)

	if err != nil {
		BadRequest(w, r, "")
		return
	}

	err = u.Create()

	if err != nil {
		BadRequest(w, r, "")
		return
	}

	render.JSON(w, r, u)
}

func userCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")
		u, err := model.UserByUsername(username)
		if err != nil {
			InternalServerError(w, r, "")
			return
		} else if (model.User{}) == u {
			NotFound(w, r, "No user found for given username.")
			return
		}

		ctx := context.WithValue(r.Context(), keyUser, u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	u := r.Context().Value(keyUser).(model.User)
	render.JSON(w, r, u)
}
