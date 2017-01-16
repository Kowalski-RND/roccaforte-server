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
