package model

import (
	"github.com/satori/go.uuid"
)

type Secret struct {
	Id         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Author     User      `json:"author"`
	CipherText string    `json:"cipher_text"`
	IV         string    `json:iv"`
}

type Secrets []Secret
