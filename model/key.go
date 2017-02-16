package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	"gopkg.in/mgutz/dat.v1/sqlx-runner"
)

// Key represents an asymmetrically encrypted
// symmetrically key. Once decrypted, this will decrypt the
// associated secret.
type Key struct {
	ID     uuid.UUID `db:"id" json:"id"`
	Secret Secret    `db:"secret" json:"secret,omitempty"`
	Owner  User      `db:"owner" json:"owner"`
	Key    string    `db:"key" json:"key"`
}

// Keys is a convenience type representing a slice of Key.
type Keys []Key

func (k Key) validateNew() error {
	return validation.StructRules{}.
		Add("Secret", validation.Required).
		Add("Owner", validation.Required).
		Add("Key", validation.Required).
		Validate(k)
}

// Create assigns a UUID and stores the Key struct
// representation into the database.
func (k Key) Create(tx *runner.Tx, s Secret) (Key, error) {
	k.ID = uuid.NewV4()
	k.Secret = s

	err := k.validateNew()

	if err != nil {
		return k, err
	}

	_, err = tx.SQL(`INSERT INTO keys (id, secret, owner, key) VALUES ($1, $2, $3, $4)`,
		&k.ID, &k.Secret.ID, &k.Owner.ID, &k.Key).
		Exec()

	return k, err
}

// AllKeysForSecret queries for associated keys. The referenced
// secret in the key is assigned to the secret passed in.
func AllKeysForSecret(secret Secret) (Keys, error) {
	keys := Keys{}

	err := db.SQL(`SELECT
				k.id, k.key, 
				u.id AS "owner.id", u.username AS "owner.username", u.fullname AS "owner.fullname"
			FROM keys k
				INNER JOIN users u ON k.owner = u.id
			WHERE k.secret = $1`, secret.ID).
		QueryStructs(&keys)

	if err != nil {
		return Keys{}, errors.Wrap(err, "Unable to query database for keys for given secret.")
	}

	for i := range keys {
		keys[i].Secret = secret
	}

	return keys, nil
}

// NukeKeysForSecret blows away all keys for a given secret
// due to the fact that they all need to be updated when a secret is updated.
func NukeKeysForSecret(secret Secret) error {
	_, err := db.DeleteFrom("keys").
		Where("secret = $1", secret.ID).
		Exec()

	return err
}
