package jwttoken

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}
