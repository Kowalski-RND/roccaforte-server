package api

import (
	"context"
	"encoding/json"
	"github.com/pressly/chi/render"
	"github.com/roccaforte/server/errors"
	"github.com/roccaforte/server/sec"
	"net/http"
)

type contextKey int

type content interface{}

type handler func(http.ResponseWriter, *http.Request) (content, error)

func (fn handler) Serve(w http.ResponseWriter, r *http.Request) {
	if c, err := fn(w, r); err != nil {
		apiError := err.(errors.Error)
		response(w, r, apiError.Error(), apiError.Status())
	} else {
		response(w, r, c, http.StatusOK)
	}
}

func response(w http.ResponseWriter, r *http.Request, c content, status int) {
	render.Status(r, status)
	if c != "" {
		render.JSON(w, r, c)
	} else {
		render.JSON(w, r, http.StatusText(status))
	}
}

func bearerTokenCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		t := r.Header.Get("Authorization")

		if t == "" {
			response(w, r, "Bearer token missing.", http.StatusBadRequest)
			return
		}

		c, err := sec.ParseJWT(t)

		if err != nil {
			response(w, r, "Invalid token provided.", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), ctxJWT, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func decode(r *http.Request, m interface{}) error {
	d := json.NewDecoder(r.Body)
	return d.Decode(m)
}
