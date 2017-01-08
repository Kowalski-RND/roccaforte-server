package model

import (
	"github.com/satori/go.uuid"
)

type Group struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type Groups []Group
