package model

import (
	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/roccaforte/server/sec"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

type User struct {
	Id        uuid.UUID `json:"-"`
	Fullname  string    `json:"fullname" validate:"required,gt=8"`
	Username  string    `json:"username"validate:"required,gt=3"`
	Password  string    `json:"password,omitempty" validate:"required,gt=3"`
	PublicKey string    `json:"public_key" validate:"required,gt=3"`
}

type Users []User

func AllUsers() (*Users, error) {
	var (
		users Users
		u     User
	)

	rows, err := db.Query(`SELECT id, username, fullname, "publicKey" FROM users`)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to query database for users.")
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.Id, &u.Username, &u.Fullname, &u.PublicKey)
		if err != nil {
			break
		}
		users = append(users, u)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return &users, err
}

func (u User) Create() error {
	u.Id = uuid.NewV4()

	hash, _ := sec.HashPwd(u.Password)

	u.Password = hash

	err := validate.Struct(u)

	if err != nil {
		return errors.Wrap(err.(validator.ValidationErrors), "Provided user data failed validation.")
	}

	_, err = db.Query(`INSERT INTO users (id, username, fullname, "publicKey", password) VALUES ($1, $2, $3, $4, $5)`,
		&u.Id, &u.Username, &u.Fullname, &u.PublicKey, &u.Password)

	return err
}

func UserByUsername(un string) (*User, bool, error) {
	var u User

	err := db.QueryRow(`SELECT id, username, password, fullname, "publicKey" FROM users WHERE username = $1`, un).
		Scan(&u.Id, &u.Username, &u.Password, &u.Fullname, &u.PublicKey)

	if err != nil && err == sql.ErrNoRows {
		return nil, true, nil
	} else if err != nil {
		return nil, false, errors.Wrap(err, "Internal Server Error while quering for user by username.")
	}

	return &u, false, nil
}

// Sets password to empty string to omit on serialization.
func (u User) MarshalJSON() ([]byte, error) {
	type Alias User
	u.Password = ""
	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(u),
	})
}
