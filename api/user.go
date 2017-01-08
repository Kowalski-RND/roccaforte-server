package api

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/roccaforte/server/model"
	"net/http"
)

func userRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", allUsers)
	r.Post("/", createUser)

	r.Route("/:username", func(r chi.Router) {
		r.Use(userCtx)
		r.Get("/", getUser)
	})

	return r
}

func allUsers(w http.ResponseWriter, r *http.Request) {
	users, err := model.AllUsers()

	if err != nil {
		InternalServerError(w, r, "")
		return
	}

	render.JSON(w, r, users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)

	var u model.User
	err := d.Decode(&u)

	if err != nil {
		BadRequest(w, r, "")
		return
	}

	defer r.Body.Close()

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
		u, empty, err := model.UserByUsername(username)
		if err != nil {
			InternalServerError(w, r, "")
			return
		} else if empty {
			NotFound(w, r, "No user found for given username.")
			return
		}
		ctx := context.WithValue(r.Context(), "user", u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUser(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").(*model.User)
	render.JSON(w, r, u)
}
