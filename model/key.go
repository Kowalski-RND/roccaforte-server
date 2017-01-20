package model

import (
	"github.com/satori/go.uuid"
)

// Key represents an asymmetrically encrypted
// symmetrically key. Once decrypted, this will decrypt the
// associated secret.
type Key struct {
	ID     uuid.UUID `json:"id"`
	Secret Secret    `json:"secret" validate:"required"`
	Owner  User      `json:"owner" validate:"required"`
	Key    string    `json:"key" validate:"required"`
}

// Keys is a convenience type representing a slice of Key.
type Keys []Key

func AllKeysForSecret(secret Secret) (Keys, error) {
	var (
		temp *Key = &Key{}
		keys Keys
	)

	rows, err := db.NamedQuery(`SELECT
			k.id, k.key, 
			u.id AS "owner.id", u.username AS "owner.username", u.fullname AS "owner.fullname"
		FROM key k
			INNER JOIN users u ON k.owner = u.id
		WHERE k.secret = :id`, secret)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&temp)

		if err != nil {
			continue
		}

		temp.Secret = secret
		keys = append(keys, *temp)
	}

	return keys, nil
}
