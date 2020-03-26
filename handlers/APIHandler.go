package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"releasetrackr/db"
	"releasetrackr/jwttoken"
	"releasetrackr/models"
	"releasetrackr/requests"
	"releasetrackr/responses"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
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

	if r.Method != "POST" {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  http.StatusMethodNotAllowed,
			Error: "Method not allowed",
		})

		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(json)

		return
	}

	decoder := json.NewDecoder(r.Body)

	var login = &requests.LoginRequest{}

	// If we couldn't decode the request, then we'll
	// send them back with a 400.
	err := decoder.Decode(login)
	if err != nil {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  400,
			Error: "Request format invalid.",
		})

		w.WriteHeader(400)
		w.Write(json)

		return
	}

	defer r.Body.Close()

	sess, _ := db.GetDbSession()

	c := sess.Database("releasetrackr").Collection("users")

	var user models.User

	res := c.FindOne(context.Background(), bson.M{
		"email": login.Email,
	})

	if err := res.Err(); err != nil {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  403,
			Error: "Username/Password invalid - USR",
		})

		w.WriteHeader(400)
		w.Write(json)

		return
	}

	res.Decode(&user)

	hashErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))

	if hashErr != nil {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  403,
			Error: "Username/Password invalid - PSW",
		})

		w.WriteHeader(400)
		w.Write(json)

		return
	}

	claims := &jwttoken.CustomClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "releasetrackr",
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signingKey := []byte(os.Getenv("JWT_SECRET"))
	tokenString, _ := token.SignedString(signingKey)

	json, _ := json.Marshal(&responses.SuccessResponse{
		Code:    200,
		Message: tokenString,
	})

	w.WriteHeader(200)
	w.Write(json)
}

// APIUserHandler provides the /api/user endpoint
func APIUserHandler(w http.ResponseWriter, r *http.Request) {
	var userModel models.User

	user, err := userModel.GetUserFromJWT(r.Header.Get("Authorization"))

	if err != nil {
		json, _ := json.Marshal(&responses.ErrorResponse{
			Code:  http.StatusInternalServerError,
			Error: "User not found",
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		return
	}

	json, _ := json.Marshal(user)
	w.Write(json)

	return
}
