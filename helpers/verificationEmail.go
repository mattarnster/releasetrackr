package helpers

import (
	"log"
	"os"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

func SendVerificationEmail(email string, vt string) {
	mg := mailgun.NewMailgun("mattarnster.co.uk", os.Getenv("MAILGUN_API_KEY"), "")
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
