package helpers

import (
	"log"
	"os"
	"time"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// SendVerificationEmail does exactly what it says on the tin.
func SendVerificationEmail(email string, vt string) {
	mg := mailgun.NewMailgun("mattarnster.co.uk", os.Getenv("MAILGUN_API_KEY"), "")
	message := mailgun.NewMessage(
		"releasetrackr@mattarnster.co.uk",
		"releasetrackr : Verify your email",
		"Hey there, please visit http://localhost:3000/verify?key="+vt+" to verify your email address so that you can receive releasetrackr notifications!",
		email)

	message.SetHtml("<h1>releasetrackr</h1><p>Hey there, please visit <a href=\"http://localhost:3000/verify?key=" + vt + "\">http://localhost:3000/verify?key=" + vt + "</a> to verify your email address.</p><p>You have until " + time.Now().Add(1*time.Hour).String() + " until your information is deleted from our systems.</p>")

	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[Helper][Mailgun] Queued ID: %s Resp: %s\n", id, resp)
}
