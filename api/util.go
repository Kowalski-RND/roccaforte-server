package api

import (
	"context"
	"github.com/roccaforte/server/sec"
	"net/http"
)

type contextKey int

func bearerTokenCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		t := r.Header.Get("Authorization")

		if t == "" {
			Unauthorized(w, r, "Bearer token missing")
			return
		}

		c, err := sec.ParseJWT(t)

		if err != nil {
			BadRequest(w, r, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), ctxJWT, c)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
