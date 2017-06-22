package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"releasetrackr/helpers"
	"releasetrackr/models"
	"releasetrackr/responses"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var token string
var user *models.User

// AuthMiddleware checks that the user is logged in
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authorizationHeader := r.Header["Authorization"]

		if len(authorizationHeader) >= 1 {
			token = authorizationHeader[0]
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  http.StatusForbidden,
				Error: "Session token missing",
			})

			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(json))
			return
		}

		if token == "" {
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  http.StatusForbidden,
				Error: "Session token missing",
			})

			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(json))
			return
		}

		db, err := helpers.GetDbSession()
		if err != nil {
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  http.StatusInternalServerError,
				Error: "Internal Server Error - E_DB_NO_CONN",
			})

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(json))

			return
		}

		sess := db.DB("releasetrackr").C("users")

		log.Printf("[AuthMiddleware] Trying to find user with session token %s", token)

		notFound := sess.Find(bson.M{
			"session_token": token,
		}).One(&user)

		if notFound != nil {
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  http.StatusForbidden,
				Error: "Session token invalid",
			})

			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(json))

			return
		}

		if user.SessionExpiry.After(time.Now()) {
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  http.StatusUnauthorized,
				Error: "Session expired, please login.",
			})

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(json))

			return
		}

		next.ServeHTTP(w, r)
	})
}
