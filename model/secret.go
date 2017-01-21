package model

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// Secret represents an symmetrically encrypted piece of data.
// It does not encompass an individual key.
type Secret struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Name       string    `db:"name" json:"name" validate:"required"`
	Author     User      `db:"author" json:"author" validate:"required"`
	Keys       Keys      `json:"keys,omitempty" validate:"required"`
	CipherText string    `db:"cipher_text" json:"cipher_text" validate:"required"`
	IV         string    `db:"iv" json:"iv" validate:"required"`
}

// Secrets is a convenience type representing a slice of Secret.
type Secrets []Secret

// AllSecrets returns all secrets authored by given user id (taken from JWT)
func AllSecrets(author uuid.UUID) (Secrets, error) {
	secrets := Secrets{}

	err := db.SQL(`SELECT 
				s.id,s.name, s.cipher_text, s.iv,
				a.fullname AS "author.fullname", a.username AS "author.username"
			FROM secrets s 
				INNER JOIN users a ON a.id = s.author
			WHERE s.author = $1`, author).
		QueryStructs(&secrets)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to query database for secrets.")
	}

	for i := range secrets {
		k, _ := AllKeysForSecret(secrets[i])
		secrets[i].Keys = append(secrets[i].Keys, k...)
	}

	return secrets, nil
}

// Create assigns a UUID and stores the Secret struct
// representation into the database.
func (s Secret) Create(author uuid.UUID) (Secret, error) {
	s.ID = uuid.NewV4()

	// Add struct validation here later

	_, err := db.DB.Query(`INSERT INTO secrets (id, name, author, cipher_text, iv) VALUES ($1, $2, $3, $4, $5)`,
		&s.ID, &s.Name, &author, &s.CipherText, &s.IV)

	return s, err
}
