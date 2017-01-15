package model

import (
	"github.com/satori/go.uuid"
)

// Secret represents an symmetrically encrypted piece of data.
// It does not encompass an individual key.
type Secret struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Author     uuid.UUID `json:"author"`
	CipherText string    `json:"cipher_text"`
	IV         string    `json:"iv"`
}

type secrets []Secret

// Create assigns a UUID and stores the Secret struct
// representation into the database.
func (s Secret) Create() (Secret, error) {
	s.ID = uuid.NewV4()

	// Add struct validation here later

	_, err := db.Query(`INSERT INTO secret (id, name, author, cipher_text, iv) VALUES ($1, $2, $3, $4, $5)`,
		&s.ID, &s.Name, &s.Author, &s.CipherText, &s.IV)

	return s, err
}
