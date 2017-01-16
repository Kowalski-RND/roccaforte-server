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
	CipherText string    `json:"cipher_text" validate:"required"`
	IV         string    `json:"iv" validate:"required"`
}

// Secrets is a convenience type representing a slice of Secret.
type Secrets []Secret

// AllSecrets returns all secrets authored by given user id (taken from JWT)
func AllSecrets(tokenUser uuid.UUID) (Secrets, error) {
	var (
		temp    map[string]Secret
		secrets Secrets
		s       Secret
		u       User
		k       Key
	)

	rows, err := db.Query(`SELECT s.id, s.name, a.id, a.username, a.fullname, s.cipher_text, s.iv, k.id, k.key, u.id, u.username, u.fullname 
		FROM secret s 
			INNER JOIN users a ON a.id = s.author
			LEFT JOIN  key   k ON s.id = k.secret 
			LEFT JOIN  users u ON k.owner = u.id 
		WHERE s.author = $1`,
		tokenUser.String())

	if err != nil {
		return nil, errors.Wrap(err, "Unable to query database for secrets.")
	}

	defer rows.Close()

	temp = make(map[string]Secret)

	for rows.Next() {
		err := rows.Scan(&s.ID, &s.Name, &s.Author.ID, &s.Author.Username, &s.Author.Fullname, &s.CipherText, &s.IV, &k.ID, &k.Key, &u.ID, &u.Username, &u.Fullname)
		if err != nil {
			break
		}

		k.Owner = u
		k.Secret = s

		if val, ok := temp[s.ID.String()]; ok {
			val.Keys = append(val.Keys, k)
		} else {
			s.Keys = append(s.Keys, k)
			temp[s.ID.String()] = s
		}
	}

	for _, val := range temp {
		secrets = append(secrets, val)
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
