package sec

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/now"
	"github.com/satori/go.uuid"
	"log"
	"time"
)

const (
	alg                = "HS256"
	iss                = "Roccaforte"
	secret             = "potato"
	invalidToken       = "Provided token is invalid."
	invalidSigningAlgo = "Signing algorithm of provided token is incorrect"
)

// IssueJWT creates a new token to be returned as part of
// a successful login response
func IssueJWT(sub uuid.UUID) (string, error) {

	log.Println("Issuing JWT")

	n := time.Now().UTC()

	c := jwt.StandardClaims{
		Id:        uuid.NewV4().String(),
		IssuedAt:  n.Unix(),
		ExpiresAt: now.New(n).EndOfWeek().Unix(),
		NotBefore: n.Unix(),
		Issuer:    iss,
		Subject:   sub.String(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	return t.SignedString([]byte(secret))
}

// ParseJWT takes a token, validates it (and the algorithm used to sign it)
// and will return a slice of claims. If something went wrong and error will be returned.
func ParseJWT(t string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(invalidSigningAlgo)
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, errors.New(invalidToken)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New(invalidToken)
	}
}
