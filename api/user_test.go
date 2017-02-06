package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/roccaforte/server/model"
	"github.com/satori/go.uuid"
)

var (
	defaultUUID = uuid.UUID{}
)

type testUser struct {
	Fullname  string `db:"fullname" json:"fullname"`
	Username  string `db:"username" json:"username"`
	Password  string `db:"password" json:"password,omitempty"`
	PublicKey string `db:"public_key" json:"public_key,omitempty"`
}

func TestUserCreation(t *testing.T) {

	u := testUser{
		Fullname:  "Brandon Kowalski",
		Username:  "BrandonKowalski",
		Password:  "ILovePotatoSalad",
		PublicKey: "This is really fake",
	}

	en, err := json.Marshal(u)

	if err != nil {
		t.Error(err.Error())
	}

	resp, err := http.Post("http://localhost:9090/users", "application/json", bytes.NewBuffer(en))

	if err != nil {
		t.Error(err.Error())
	}

	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	c := model.User{}

	err = d.Decode(&c)

	if err != nil {
		t.Error(err.Error())
	}

	if c.ID == defaultUUID {
		t.Error("User not saved as default UUID present on returned data")
	}
}
