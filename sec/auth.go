package sec

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPwd returns the bcrypt hash of the supplied password.
// If something goes wrong (it shouldn't) this returns empty string along with an error.
func HashPwd(p string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(h), nil
}

// CheckPw compares a bcrypt hash password against a clear text password.
func CheckPw(h, p string) bool {
	return bcrypt.CompareHashAndPassword([]byte(h), []byte(p)) == nil
}
