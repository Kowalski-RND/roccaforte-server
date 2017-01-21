package api

import (
	"encoding/json"
	"github.com/pressly/chi"
	"github.com/pressly/chi/render"
	"github.com/roccaforte/server/model"
	"github.com/roccaforte/server/sec"
	"net/http"
)

const (
	invalidUserOrPass string = "Invalid username or password."
)

func loginRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", login)

	return r
}

func login(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	d := json.NewDecoder(r.Body)

	var c model.Credentials
	err := d.Decode(&c)

	if err != nil {
		BadRequest(w, r, "")
		return
	}

	u, err := model.UserByUsername(c.Username)

	if err != nil {
		InternalServerError(w, r, "")
		return
	} else if (model.User{}) == u {
		Unauthorized(w, r, invalidUserOrPass)
		return
	}

	a := sec.CheckPw(u.Password, c.Password)

	if !a {
		Unauthorized(w, r, invalidUserOrPass)
		return
	}

	t, err := sec.IssueJWT(u.ID)

	if err != nil {
		InternalServerError(w, r, "")
		return
	}

	render.JSON(w, r, t)
}
