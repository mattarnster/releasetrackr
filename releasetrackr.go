package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/mattarnster/releasetrackr/handlers"
	"github.com/mattarnster/releasetrackr/jobs"
	"github.com/mattarnster/releasetrackr/middleware"
)

func main() {
	log.Println("[App][Startup] releasetrackr - 1.0 started")

	if os.Getenv("MAILGUN_API_KEY") != "" {
		log.Println("Mailgun API Key detected.")
	} else {
		panic("Couldn't get Mailgun API key from environment variable MAILGUN_API_KEY, make sure this is set.")
	}

	if os.Getenv("MONGO_HOST") != "" {
		log.Println("Using MongoDB Host: " + os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT"))
	} else {
		panic("Environment variable doesn't exist or is empty: MONGO_HOST - Please make sure it is present and correct.")
	}

	// HTTP Handlers
	httpIndex := http.HandlerFunc(handlers.IndexHandler)
	httpTrack := http.HandlerFunc(handlers.TrackHandler)
	httpVerify := http.HandlerFunc(handlers.VerificationHandler)

	http.Handle("/", middleware.ContentTypeMiddleware(httpIndex))
	http.Handle("/track", middleware.ContentTypeMiddleware(httpTrack))
	http.Handle("/verify", middleware.ContentTypeMiddleware(httpVerify))

	// Setting up scheduled jobs
	go func() {
		gocron.Every(1).Hour().Do(jobs.ClearNonVerifiedUsers)
		gocron.Every(2).Hours().Do(jobs.GetNewReleases)
		gocron.RunAll()
		<-gocron.Start()
	}()

	http.ListenAndServe(":3000", nil)
}
