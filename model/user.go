package model

import (
	"database/sql"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/roccaforte/server/sec"
	"github.com/satori/go.uuid"
	"gopkg.in/go-playground/validator.v9"
)

// User represents an individual that has access to Roccaforte
type User struct {
	ID        uuid.UUID `json:"-"`
	Fullname  string    `json:"fullname" validate:"required,gt=8"`
	Username  string    `json:"username"validate:"required,gt=3"`
	Password  string    `json:"password,omitempty" validate:"required,gt=3"`
	PublicKey string    `json:"public_key" validate:"required,gt=3"`
}

type users []User

// Credentials is a convenience struct for login use.
type Credentials struct {
	Username string `json:"username`
	Password string `json:"password"`
}

// AllUsers retreives all users from the database.
func AllUsers() (*users, error) {
	var (
		users users
		u     User
	)

	rows, err := db.Query(`SELECT id, username, fullname, "publicKey" FROM users`)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to query database for users.")
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Fullname, &u.PublicKey)
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

// Create assigns a UUID and stores the User struct
// representation into the database.
func (u User) Create() error {
	u.ID = uuid.NewV4()

	hash, _ := sec.HashPwd(u.Password)

	u.Password = hash

	err := validate.Struct(u)

	if err != nil {
		return errors.Wrap(err.(validator.ValidationErrors), "Provided user data failed validation.")
	}

	_, err = db.Query(`INSERT INTO users (id, username, fullname, "publicKey", password) VALUES ($1, $2, $3, $4, $5)`,
		&u.ID, &u.Username, &u.Fullname, &u.PublicKey, &u.Password)

	return err
}

// UserByUsername retreives a user based on provided username.
// If no user exists for given username, middle return will be true.
// If an error occurs during the query, middle return will be false and error will be populated.
func UserByUsername(un string) (*User, bool, error) {
	var u User

	err := db.QueryRow(`SELECT id, username, password, fullname, "publicKey" FROM users WHERE username = $1`, un).
		Scan(&u.ID, &u.Username, &u.Password, &u.Fullname, &u.PublicKey)

	if err != nil && err == sql.ErrNoRows {
		return nil, true, nil
	} else if err != nil {
		return nil, false, errors.Wrap(err, "Internal Server Error while quering for user by username.")
	}

	return &u, false, nil
}

// MarshalJSON overrides default functionality and sets
// password to empty string to omit on serialization.
func (u User) MarshalJSON() ([]byte, error) {
	type alias User
	u.Password = ""
	return json.Marshal(&struct {
		alias
	}{
		alias: (alias)(u),
	})
}
