package model

import (
	"database/sql"
	"encoding/json"

	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/roccaforte/server/sec"
	"github.com/satori/go.uuid"
)

// User represents an individual that has access to Roccaforte
type User struct {
	ID        uuid.UUID `db:"id" json:"id,omitempty"`
	Fullname  string    `db:"fullname" json:"fullname"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password,omitempty"`
	PublicKey string    `db:"public_key" json:"public_key,omitempty"`
}

// Users is a convenience type representing a slice of User.
type Users []User

// Credentials is a convenience struct for login use.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate ensures that nested usage of User contains an ID
func (u User) Validate() error {
	return validation.StructRules{}.
		Add("ID", validation.Required).
		Validate(u)
}

func (u User) validateNew() error {
	return validation.StructRules{}.
		Add("Fullname", validation.Required).
		Add("Username", validation.Required, validation.Length(3, 30)).
		Add("Password", validation.Required, validation.Length(8, 256)).
		// Add custom validator for Public Key verification
		Validate(u)
}

// AllUsers retreives all users from the database.
func AllUsers() (Users, error) {
	users := Users{}

	err := db.
		Select("id", "username", "fullname", "public_key").
		From("users").
		QueryStructs(&users)

	if err != nil {
		return Users{}, errors.Wrap(err, "Unable to query database for users.")
	}

	return users, nil
}

// Create assigns a UUID and stores the User struct
// representation into the database.
func (u User) Create() (User, error) {
	u.ID = uuid.NewV4()

	err := u.validateNew()

	if err != nil {
		return u, err
	}

	u.Password, err = sec.HashPwd(u.Password)

	if err != nil {
		return u, err
	}

	err = db.
		InsertInto("users").
		Columns("id", "username", "fullname", "public_key", "password").
		Values(u.ID, u.Username, u.Fullname, u.PublicKey, u.Password).
		Returning("id").
		QueryStruct(&u)

	return u, err
}

// UserByUsername retreives a user based on provided username.
func UserByUsername(un string) (User, error) {
	u := User{}

	err := db.
		Select("id", "username", "password", "fullname", "public_key").
		From("users").
		Where("username = $1", un).
		QueryStruct(&u)

	if err != nil && err != sql.ErrNoRows {
		return User{}, errors.Wrap(err, "Internal Server Error while quering for user by username.")
	}

	return u, nil
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
