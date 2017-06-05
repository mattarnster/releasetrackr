package helpers

import (
	"log"
	"os"

	"github.com/mattarnster/releasetrackr/models"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

// SendNotificationEmail does exactly what it says on the tin.
func SendNotificationEmail(repo models.Repo, email string, release models.Release) {
	mg := mailgun.NewMailgun("mattarnster.co.uk", os.Getenv("MAILGUN_API_KEY"), "")
	message := mailgun.NewMessage(
		"releasetrackr@mattarnster.co.uk",
		"releasetrackr : New Release for "+repo.Repo,
		"Hey there, A new release for "+repo.Repo+" has been detected!",
		email)

	message.SetHtml("<h1>releasetrackr</h1><p>Hey there.</p><p>A new release for " + repo.Repo + " has been detected! " + release.Body)

	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[Helper][NotificationEmail][Mailgun] Queued ID: %s Resp: %s\n", id, resp)
}
