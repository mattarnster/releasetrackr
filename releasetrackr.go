package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"

	"releasetrackr/handlers"
	"releasetrackr/jobs"
	"releasetrackr/middleware"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found or error loading it. Using environment variables.")
	} else {
		log.Println("[Startup] Loaded environment variables from .env file")
	}

	log.Println("[Startup] releasetrackr - 2.0 started")

	if os.Getenv("MAILJET_API_PUBLIC_KEY") != "" {
		log.Println("[Startup] Mailjet API Public Key detected.")
	} else {
		panic("Couldn't get Mailjet API Public Key from environment variable MAILJET_API_PUBLIC_KEY, make sure this is set.")
	}

	if os.Getenv("MAILJET_API_PRIVATE_KEY") != "" {
		log.Println("[Startup] Mailjet API Private Key detected.")
	} else {
		panic("Couldn't get Mailjet API Private Key from environment variable MAILJET_API_PRIVATE_KEY, make sure this is set.")
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

	if os.Getenv("RECAPTCHA_SECRET") != "" {
		log.Println("[Startup] RECAPTCHA_SECRET detected.")
	} else {
		panic("Didn't find RECAPTCHA_SECRET in environment, please set it.")
	}

	if os.Getenv("JWT_SECRET") != "" {
		log.Println("[Startup] JWT_SECRET detected.")
	} else {
		panic("Didn't find JWT_SECRET in environment, please set it.")
	}

	// HTTP Handlers
	httpIndex := http.HandlerFunc(handlers.IndexHandler)
	httpTrack := http.HandlerFunc(handlers.TrackHandler)
	httpVerify := http.HandlerFunc(handlers.VerificationHandler)
	httpStats := http.HandlerFunc(handlers.StatsHandler)
	apiIndex := http.HandlerFunc(handlers.APIIndexHandler)
	apiToken := http.HandlerFunc(handlers.APITokenHandler)
	apiUser := http.HandlerFunc(handlers.APIUserHandler)
	// Web
	http.Handle("/", httpIndex)
	http.Handle("/track", middleware.ContentTypeMiddleware(httpTrack))
	http.Handle("/verify", middleware.ContentTypeMiddleware(httpVerify))
	http.Handle("/stats", middleware.ContentTypeMiddleware(httpStats))

	// API
	http.Handle("/api", middleware.ContentTypeMiddleware(apiIndex))
	http.Handle("/api/auth", middleware.ContentTypeMiddleware(apiToken))
	http.Handle("/api/user", middleware.ContentTypeMiddleware(middleware.AuthMiddleware(apiUser)))

	// Assets for the email templates
	fs := http.FileServer(http.Dir("assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	// Setting up scheduled jobs
	go func() {
		// gocron.Every(1).Hour().Do(jobs.ClearNonVerifiedUsers)
		// //gocron.Every(2).Hours().Do(jobs.GetNewReleases)
		// gocron.RunAll()
		// <-gocron.Start()
		for true {
			fmt.Println("Running scheduled job: ClearNonVerifiedUsers")
			jobs.ClearNonVerifiedUsers()
			time.Sleep(2 * time.Hour)
		}
	}()

	go func() {
		// gocron.Every(1).Hour().Do(jobs.ClearNonVerifiedUsers)
		// //gocron.Every(2).Hours().Do(jobs.GetNewReleases)
		// gocron.RunAll()
		// <-gocron.Start()
		for true {
			fmt.Println("Running scheduled job: GetNewReleases")
			jobs.GetNewReleases()
			time.Sleep(5 * time.Minute)
		}
	}()

	http.ListenAndServe(":3000", nil)
}
