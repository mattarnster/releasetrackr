package jwttoken

import (
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// GetJWT for when you need to pull a
// token out of the header.
func GetJWT(header string) (*jwt.Token, error) {
	jwtTok := strings.Split(header, "Bearer ")

	tok, err := jwt.ParseWithClaims(jwtTok[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	return tok, nil
}
