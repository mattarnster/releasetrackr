package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mattarnster/releasetrackr/handlers"
	"github.com/mattarnster/releasetrackr/middleware"
)

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

	http.Handle("/", middleware.ContentTypeMiddleware(httpIndex))
	http.Handle("/track", middleware.ContentTypeMiddleware(httpTrack))
	http.Handle("/verify", middleware.ContentTypeMiddleware(httpVerify))
	http.ListenAndServe(":3000", nil)
}
