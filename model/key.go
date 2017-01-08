package model

import (
	"github.com/satori/go.uuid"
)

type Key struct {
	Id     uuid.UUID `json:"id"`
	Secret Secret    `json:"secret"`
	Owner  User      `json:"user"`
	Key    string    `json:"key"`
}

type Keys []Key
