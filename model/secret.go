package model

import (
	"github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	runner "gopkg.in/mgutz/dat.v1/sqlx-runner"
)

// Secret represents an symmetrically encrypted piece of data.
// It does not encompass an individual key.
type Secret struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Author     User      `db:"author" json:"author"`
	Keys       Keys      `json:"keys"`
	CipherText string    `db:"cipher_text" json:"cipher_text"`
	IV         string    `db:"iv" json:"iv"`
}

// Secrets is a convenience type representing a slice of Secret.
type Secrets []Secret

// Validate ensures that nested usage of Secret contains an ID
func (s Secret) Validate() error {
	return validation.StructRules{}.
		Add("ID", validation.Required).
		Validate(s)
}

func (s Secret) validateNew() error {
	return validation.StructRules{}.
		Add("Author", validation.Required).
		Add("Keys", validation.Required).
		Add("CipherText", validation.Required).
		Add("IV", validation.Required).
		Validate(s)
}

// AllSecrets returns all secrets authored by given user id (taken from JWT)
func AllSecrets(author uuid.UUID) (Secrets, error) {
	secrets := Secrets{}

	err := db.SQL(`SELECT 
				s.id, s.cipher_text, s.iv,
				a.fullname AS "author.fullname", a.username AS "author.username"
			FROM secrets s 
				INNER JOIN users a ON a.id = s.author
			WHERE s.author = $1`, author).
		QueryStructs(&secrets)

	if err != nil {
		return Secrets{}, errors.Wrap(err, "Unable to query database for secrets.")
	}

	for i := range secrets {
		k, _ := AllKeysForSecret(secrets[i])
		secrets[i].Keys = append(secrets[i].Keys, k...)
	}

	return secrets, nil
}

// GetSecret retreives a Secret for the given ID.
func GetSecret(id uuid.UUID) (Secret, error) {
	s := Secret{}

	err := db.SQL(`SELECT 
				s.id, s.cipher_text, s.iv,
				a.id AS "author.id", a.fullname AS "author.fullname", a.username AS "author.username"
			FROM secrets s 
				INNER JOIN users a ON a.id = s.author
			WHERE s.id = $1`, id).
		QueryStruct(&s)

	return s, err
}

// Create assigns a UUID and stores the Secret struct
// representation into the database.
func (s Secret) Create() (Secret, error) {
	tx, _ := CreateTransaction()
	defer tx.AutoRollback()

	s.ID = uuid.NewV4()

	err := s.validateNew()

	if err != nil {
		return s, err
	}

	_, err = tx.SQL(`INSERT INTO secrets (id, author, cipher_text, iv) VALUES ($1, $2, $3, $4)`,
		&s.ID, &s.Author.ID, &s.CipherText, &s.IV).
		Exec()

	err = insertKeys(tx, s)

	if err != nil {
		return s, errors.Wrap(err, "Unable to save keys.")
	}

	tx.Commit()

	return s, err
}

func insertKeys(tx *runner.Tx, s Secret) error {
	for i := range s.Keys {
		s.Keys[i].Secret = s
		k, err := s.Keys[i].Create(tx, s)

		if err != nil {
			return err
		}

		// Prevent circular, stack destroying struct.
		s.Keys[i].Secret.Keys = Keys{}
		s.Keys[i].ID = k.ID
	}

	return nil
}

// Update modifies a saved Secret. This will delete all current Keys
// associated with the Secret and replace all fields with what got provided.
func (s Secret) Update() (Secret, error) {
	tx, _ := CreateTransaction()
	defer tx.AutoRollback()

	err := s.Validate()

	if err != nil {
		return s, errors.Wrap(err, "Cannot update Secret as no ID passed in.")
	}

	_, err = tx.
		Update("secrets").
		Set("cipher_text", s.CipherText).
		Set("iv", s.IV).
		Where("id = $1", s.ID).
		Exec()

	if err != nil {
		return s, errors.Wrap(err, "Unable to update secret. No changes have been persisted.")
	}

	err = NukeKeysForSecret(s)

	if err != nil {
		return s, errors.Wrap(err, "Unable to delete existing keys.")
	}

	err = insertKeys(tx, s)

	if err != nil {
		return s, errors.Wrap(err, "Unable to save keys.")
	}

	tx.Commit()

	return s, err
}

// DeleteSecret removes a saved Secret. Please note to any consumer of this function
// that this removes the Secret and associated Keys from the database. You still
// have to roll the password in this Secret as individuals may still hold a copy.
// Common sense, but just a friendly reminder :)
func (s Secret) Delete() error {
	tx, _ := CreateTransaction()
	defer tx.AutoRollback()

	err := s.Validate()

	if err != nil {
		return errors.Wrap(err, "Cannot delete Secret as no ID passed in.")
	}

	_, err = db.DeleteFrom("secrets").
		Where("id = $1", s.ID).
		Exec()

	return err
}
