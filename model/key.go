package model

import (
	"github.com/satori/go.uuid"
)

// Key represents an asymmetrically encrypted
// symmetrically key. Once decrypted, this will decrypt the
// associated secret.
type Key struct {
	ID     uuid.UUID `json:"id"`
	Secret Secret    `json:"secret"`
	Owner  User      `json:"owner"`
	Key    string    `json:"key"`
}

type keys []Key
