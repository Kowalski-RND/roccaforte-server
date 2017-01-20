package model

import (
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

// Secret represents an symmetrically encrypted piece of data.
// It does not encompass an individual key.
type Secret struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name" validate:"required"`
	Author     User      `json:"author" validate:"required"`
	Keys       Keys      `json:"keys,omitempty" validate:"required"`
	CipherText string    `json:"cipher_text" validate:"required" db:"cipher_text"`
	IV         string    `json:"iv" validate:"required"`
}

// Secrets is a convenience type representing a slice of Secret.
type Secrets []Secret

// AllSecrets returns all secrets authored by given user id (taken from JWT)
func AllSecrets(tokenUser uuid.UUID) (Secrets, error) {
	var (
		temp    *Secret = &Secret{}
		secrets Secrets
		author  = User{
			ID: tokenUser,
		}
	)

	rows, err := db.NamedQuery(`SELECT 
			s.id,s.name, s.cipher_text, s.iv,
			a.fullname AS "author.fullname", a.username AS "author.username"
		FROM secret s 
			INNER JOIN users a ON a.id = s.author
		WHERE s.author = :id`, author)

	if err != nil {
		return nil, errors.Wrap(err, "Unable to query database for secrets.")
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.StructScan(&temp)

		if err != nil {
			continue
		}

		secrets = append(secrets, *temp)

		keys, err := AllKeysForSecret(*temp)

		if err != nil {
			continue
		}

		temp.Keys = keys
	}

	return secrets, nil
}

// Create assigns a UUID and stores the Secret struct
// representation into the database.
func (s Secret) Create(author uuid.UUID) (Secret, error) {
	s.ID = uuid.NewV4()

	// Add struct validation here later

	_, err := db.Query(`INSERT INTO secret (id, name, author, cipher_text, iv) VALUES ($1, $2, $3, $4, $5)`,
		&s.ID, &s.Name, &author, &s.CipherText, &s.IV)

	return s, err
}
