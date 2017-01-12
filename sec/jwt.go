package sec

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/now"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

const (
	alg    = "HS256"
	iss    = "Roccaforte"
	secret = "potato"
)

// IssueJWT creates a new token to be returned as part of
// a successful login response
func IssueJWT(sub string) (string, error) {

	log.Println("Issuing JWT")

	n := time.Now().UTC()

	c := jwt.StandardClaims{
		Id:        uuid.NewV4().String(),
		IssuedAt:  n.Unix(),
		ExpiresAt: now.New(n).EndOfWeek().Unix(),
		NotBefore: n.Unix(),
		Issuer:    iss,
		Subject:   sub,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return t.SignedString([]byte(secret))
}
