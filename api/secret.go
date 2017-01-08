package api

import (
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/roccaforte/server/sec"
	"net/http"
)

func secretRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", allSecrets)

	return r
}

func allSecrets(w http.ResponseWriter, r *http.Request) {
	t, err := sec.IssueJWT()

	if err != nil {
		InternalServerError(w, r, "")
		return
	}

	render.JSON(w, r, t)
}
