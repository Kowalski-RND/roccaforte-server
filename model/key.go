package model

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// Key represents an asymmetrically encrypted
// symmetrically key. Once decrypted, this will decrypt the
// associated secret.
type Key struct {
	ID     uuid.UUID `db:"id" json:"id"`
	Secret Secret    `db:"secret" json:"secret" validate:"required"`
	Owner  User      `db:"owner" json:"owner" validate:"required"`
	Key    string    `db:"key" json:"key" validate:"required"`
}

// Keys is a convenience type representing a slice of Key.
type Keys []Key

func (k Key) CreateKey(s Secret) (Key, error) {
	k.ID = uuid.NewV4()
	k.Secret = s

	_, err := db.DB.Query(`INSERT INTO keys (id, secret, owner, key) VALUES ($1, $2, $3, $4)`,
		&k.ID, &k.Secret.ID, &k.Owner.ID, &k.Key)

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
