package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	uuid "github.com/nu7hatch/gouuid"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DbConnection Will hold the global DB Connection
var DbConnection *mgo.Session

// mailgunApiKey is set from environment variable MAILGUN_API_KEY
var mailgunAPIKey string

func getDbSession() (*mgo.Session, error) {
	if DbConnection == nil {
		session, err := mgo.Dial("localhost:32768")
		if err != nil {
			return nil, err
		}
		DbConnection = session
	}
	return DbConnection, nil
}

// KillDbSession kills the currently active DB session
func KillDbSession() {
	if DbConnection != nil {
		DbConnection.Close()
	}
}

type indexResponse struct {
	Name string `json:"name"`
	Ver  string `json:"version"`
}

type trackRequest struct {
	Repo  string `json:"repo"`
	Email string `json:"email"`
}

type errorResponse struct {
	Code  int    `json:"code"`
	Error string `json:"error"`
}

type successResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type userRecord struct {
	ID               bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Email            string        `json:"email" bson:"email,omitempty" valid:"email"`
	VerificationCode string        `json:"code" bson:"verificationcode"`
	Verified         bool          `json:"verified" bson:"verified"`
}

type trackRecord struct {
	ID     bson.ObjectId `json:"id" bson:"_id,omitempty"`
	UserID bson.ObjectId `json:"userID" bson:"userID,omitempty"`
	Repo   string        `json:"repo" bson:"repo,omitempty"`
}

func sendVerificationEmail(email string, vt string) {
	mg := mailgun.NewMailgun("mattarnster.co.uk", mailgunAPIKey, "")
	message := mailgun.NewMessage(
		"releasetrackr@mattarnster.co.uk",
		"releasetrackr : Verify your email",
		"Hey there, please visit http://localhost:3000/verify?key="+vt+" to verify your email address so that you can receive releasetrackr notifications!",
		email)
	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("ID: %s Resp: %s\n", id, resp)
}

func index(w http.ResponseWriter, r *http.Request) {
	json, _ := json.Marshal(&indexResponse{
		Name: "releasetrackr",
		Ver:  "1.0",
	})
	w.Write(json)
}

func verify(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")

	if key != "" {
		sess, err := getDbSession()
		if err != nil {
			panic("test")
		}

		user := &userRecord{}
		c := sess.DB("releasetrackr").C("users")

		userErr := c.Find(bson.M{"verificationcode": key, "verified": false}).One(&user)
		if userErr != nil {
			json, _ := json.Marshal(&errorResponse{
				Code:  400,
				Error: "The token you want to verify is invalid",
			})
			w.WriteHeader(400)
			w.Write(json)
			log.Printf("Verification token fail: %s", r.RemoteAddr)
			return
		}

		change := bson.M{"$set": bson.M{"verified": true}}
		c.Update(user, change)

		log.Printf("Verification token pass: %s - %s", key, r.RemoteAddr)

		json, _ := json.Marshal(&successResponse{
			Code:    200,
			Message: "Verification passed.",
		})

		w.WriteHeader(200)
		w.Write(json)
		return
	}
}

func track(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		json, _ := json.Marshal(&errorResponse{
			Code:  400,
			Error: "Bad Request",
		})
		w.WriteHeader(400)
		w.Write(json)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var tr = &trackRequest{}

	err := decoder.Decode(tr)
	if err != nil {
		json, _ := json.Marshal(&errorResponse{
			Code:  400,
			Error: "Request format invalid.",
		})

		w.WriteHeader(400)
		w.Write(json)

		return
	}

	defer r.Body.Close()

	log.Printf("Incoming track request: %s from %s", tr.Repo, r.RemoteAddr)

	sess, err := getDbSession()
	if err != nil {
		panic("Couldn't get DB session")
	}

	user := &userRecord{}

	c := sess.DB("releasetrackr").C("users")

	userErr := c.Find(bson.M{"email": tr.Email}).One(&user)
	if userErr != nil {

		uid := bson.NewObjectId()
		verification, _ := uuid.NewV4()

		c.Insert(&userRecord{
			ID:               uid,
			Email:            tr.Email,
			VerificationCode: verification.String(),
			Verified:         false,
		})

		log.Printf("New user: %s, %s", uid, tr.Email)
		sendVerificationEmail(tr.Email, verification.String())

		json, _ := json.Marshal(&successResponse{
			Code:    202,
			Message: "Email verification required",
		})

		w.WriteHeader(202)
		w.Write(json)

		return

	}
	// Existing user, make sure they're verified first...
	if user.Verified == false {
		response, _ := json.Marshal(&errorResponse{
			Code:  403,
			Error: "Verification required - Check your email.",
		})
		log.Println("Responding with verification required.")
		w.WriteHeader(403)
		w.Write(response)
		return
	}
	// Already a user, stop them from making another record of the same.
	c = sess.DB("releasetrackr").C("tracks")
	record := &trackRecord{}
	dbtr := c.Find(bson.M{"userID": bson.ObjectId(user.ID), "repo": tr.Repo}).One(&record)
	if dbtr == nil {
		response, _ := json.Marshal(&errorResponse{
			Code:  409,
			Error: "You've already subscribed to be notified about this repository.",
		})

		log.Printf("User already subscribed to repo: %s - %s", user.Email, tr.Repo)

		w.WriteHeader(409)
		w.Write(response)
		return
	}

	trID := bson.NewObjectId()
	c.Insert(&trackRecord{
		ID:     trID,
		UserID: user.ID,
		Repo:   tr.Repo,
	})

	log.Printf("New track request: %s from %s", trID.String(), user.Email)

	json, _ := json.Marshal(&successResponse{
		Code:    201,
		Message: "Track request acknowledged.",
	})
	w.WriteHeader(201)
	w.Write(json)

	return
}

func responseFormatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("releasetrackr - 0.1 started")

	if os.Getenv("MAILGUN_API_KEY") != "" {
		mailgunAPIKey = os.Getenv("MAILGUN_API_KEY")
		log.Println("Mailgun API Key detected.")
	} else {
		panic("Couldn't get Mailgun API key from environment variable MAILGUN_API_KEY, make sure this is set.")
	}

	httpIndex := http.HandlerFunc(index)
	httpTrack := http.HandlerFunc(track)
	httpVerify := http.HandlerFunc(verify)

	http.Handle("/", responseFormatMiddleware(httpIndex))
	http.Handle("/track", responseFormatMiddleware(httpTrack))
	http.Handle("/verify", responseFormatMiddleware(httpVerify))
	http.ListenAndServe(":3000", nil)
}
