package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jasonlvhit/gocron"
	"releasetrackr/handlers"
	"releasetrackr/jobs"
	"releasetrackr/middleware"
)

func main() {
	log.Println("[App][Startup] releasetrackr - 1.0 started")

	if os.Getenv("MAILGUN_API_KEY") != "" {
		log.Println("[Startup] Mailgun API Key detected.")
	} else {
		panic("Couldn't get Mailgun API key from environment variable MAILGUN_API_KEY, make sure this is set.")
	}

	if os.Getenv("MONGO_HOST") != "" {
		log.Println("[Startup] Using MongoDB Host: " + os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT"))
	} else {
		panic("Environment variable doesn't exist or is empty: MONGO_HOST - Please make sure it is present and correct.")
	}

	if os.Getenv("RT_DOMAIN") != "" {
		log.Println("[Startup] RT_DOMAIN is " + os.Getenv("RT_DOMAIN"))
	} else {
		panic("Didn't find RT_DOMAIN in environment, please set it so I know where I am.")
	}

	// HTTP Handlers
	httpIndex := http.HandlerFunc(handlers.IndexHandler)
	httpTrack := http.HandlerFunc(handlers.TrackHandler)
	httpVerify := http.HandlerFunc(handlers.VerificationHandler)
	httpStats := http.HandlerFunc(handlers.StatsHandler)

	http.Handle("/", httpIndex)
	http.Handle("/track", middleware.ContentTypeMiddleware(httpTrack))
	http.Handle("/verify", middleware.ContentTypeMiddleware(httpVerify))
	http.Handle("/stats", middleware.ContentTypeMiddleware(httpStats))

	// Setting up scheduled jobs
	go func() {
		gocron.Every(1).Hour().Do(jobs.ClearNonVerifiedUsers)
		gocron.Every(2).Hours().Do(jobs.GetNewReleases)
		gocron.RunAll()
		<-gocron.Start()
	}()

	http.ListenAndServe(":3000", nil)
}
