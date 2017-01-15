package model

import (
	"github.com/satori/go.uuid"
)

type Secret struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Author     uuid.UUID `json:"author"`
	CipherText string    `json:"cipher_text"`
	IV         string    `json:"iv"`
}

type Secrets []Secret

func (s Secret) Create() (Secret, error) {
	s.ID = uuid.NewV4()

	// Add struct validation here later

	_, err := db.Query(`INSERT INTO secret (id, name, author, cipher_text, iv) VALUES ($1, $2, $3, $4, $5)`,
		&s.ID, &s.Name, &s.Author, &s.CipherText, &s.IV)

	return s, err
}
