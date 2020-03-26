package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"releasetrackr/jwttoken"
	"releasetrackr/responses"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware intercepts JWT tokens and validates them
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtHeader := r.Header.Get("Authorization")
		jwtTok := strings.Split(jwtHeader, "Bearer ")

		if len(jwtTok) != 2 {
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  403,
				Error: "Missing token - check Authorization header",
			})

			w.WriteHeader(400)
			w.Write(json)

			return
		}

		//log.Printf("[Middleware][AuthMiddleware] Got JWT: %v\n", jwtTok[1])

		// First parse
		tok, err := jwt.ParseWithClaims(jwtTok[1], &jwttoken.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			log.Printf("[Middleware][AuthMiddleware] JWT error: %v\n", err.Error())
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  403,
				Error: "Token error",
			})

			w.WriteHeader(400)
			w.Write(json)

			return
		}

		// Any issues with claims?
		claims, ok := tok.Claims.(*jwttoken.CustomClaims)
		if !ok {
			log.Printf("[Middleware][AuthMiddleware] Claims error")
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  403,
				Error: "Token error",
			})

			w.WriteHeader(400)
			w.Write(json)

			return
		}

		// Has it expired?
		if claims.ExpiresAt < time.Now().Unix() {
			log.Printf("[Middleware][AuthMiddleware] Token expired")
			json, _ := json.Marshal(&responses.ErrorResponse{
				Code:  403,
				Error: "Token expired",
			})

			w.WriteHeader(400)
			w.Write(json)

			return
		}

		next.ServeHTTP(w, r)
	})
}
