package model

import (
	v "gopkg.in/go-playground/validator.v9"
)

var validate *v.Validate

func init() {
	validate = v.New()
}
