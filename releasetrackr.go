package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mattarnster/releasetrackr/handlers"
)

func responseFormatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("releasetrackr - 0.1 started")

	if os.Getenv("MAILGUN_API_KEY") != "" {
		log.Println("Mailgun API Key detected.")
	} else {
		panic("Couldn't get Mailgun API key from environment variable MAILGUN_API_KEY, make sure this is set.")
	}

	httpIndex := http.HandlerFunc(handlers.IndexHandler)
	httpTrack := http.HandlerFunc(handlers.TrackHandler)
	httpVerify := http.HandlerFunc(handlers.VerificationHandler)

	http.Handle("/", responseFormatMiddleware(httpIndex))
	http.Handle("/track", responseFormatMiddleware(httpTrack))
	http.Handle("/verify", responseFormatMiddleware(httpVerify))
	http.ListenAndServe(":3000", nil)
}
