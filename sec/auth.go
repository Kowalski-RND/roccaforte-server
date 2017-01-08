package sec

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPwd(p string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(h), nil
}
