package api

import (
	"net/http"

	"github.com/pressly/chi"
	"github.com/roccaforte/server/errors"
	"github.com/roccaforte/server/model"
	"github.com/roccaforte/server/sec"
)

const (
	invalidUserOrPass string = "Invalid username or password."
)

func loginRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", handler(login).Serve)

	return r
}

func login(w http.ResponseWriter, r *http.Request) (content, error) {
	c := model.Credentials{}
	err := decode(r, &c)

	if err != nil {
		return nil, errors.BadRequest("")
	}

	u, err := model.UserByUsername(c.Username)

	if err != nil {
		return nil, errors.InternalServerError("")
	} else if (model.User{}) == u {
		return nil, errors.Unauthorized(invalidUserOrPass)
	}

	a := sec.CheckPw(u.Password, c.Password)

	if !a {
		return nil, errors.Unauthorized(invalidUserOrPass)
	}

	t, err := sec.IssueJWT(u.ID)

	if err != nil {
		return nil, errors.InternalServerError("")
	}

	return model.Token{t}, nil
}
