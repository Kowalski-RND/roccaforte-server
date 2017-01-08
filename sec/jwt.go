package sec

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	secret = "potato"
)

func IssueJWT() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "btk44",
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	return token.SignedString([]byte(secret))
}
