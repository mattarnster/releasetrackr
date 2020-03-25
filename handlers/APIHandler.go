package handlers

import (
	"encoding/json"
	"net/http"
	"releasetrackr/responses"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// APIIndexHandler endpoint is for information about the API
func APIIndexHandler(w http.ResponseWriter, r *http.Request) {
	resp, _ := json.Marshal(&responses.IndexResponse{
		Name: "releasetrackr",
		Ver:  "1.0",
	})

	w.Write(resp)
}

// APITokenHandler endpoint will provide the user a signed token
func APITokenHandler(w http.ResponseWriter, r *http.Request) {
	signingKey := []byte("thisisasigningkey")
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = "mattarnster.co.uk@gmail.com"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(signingKey)
	w.Write([]byte(tokenString))
}
