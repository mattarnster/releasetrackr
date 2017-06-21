package helpers

import (
	"bytes"
	"html/template"
	"log"
	"os"

	"releasetrackr/models"

	"github.com/russross/blackfriday"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type templateVars struct {
	RepoName    string
	ReleaseBody template.HTML
	ReleaseTag  string
	RTDomain    string
}

// SendNotificationEmail does exactly what it says on the tin.
func SendNotificationEmail(repo models.Repo, email string, release models.Release) {
	mg := mailgun.NewMailgun("mattarnster.co.uk", os.Getenv("MAILGUN_API_KEY"), "")
	message := mailgun.NewMessage(
		"releasetrackr <releasetrackr@mattarnster.co.uk>",
		"New Release for "+repo.Repo,
		"Hey there, A new release for "+repo.Repo+" has been detected!",
		email)

	t, _ := template.ParseFiles("templates/new-release.html")
	var doc bytes.Buffer
	err := t.Execute(&doc, generateTemplateVars(repo, release))

	if err != nil {
		log.Printf("Template parse failed: %v", err.Error())
	}

	message.SetHtml(doc.String())

	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[Helper][NotificationEmail][Mailgun] Queued ID: %s Resp: %s\n", id, resp)
}

// This function returns a struct of the required template variables.
func generateTemplateVars(repo models.Repo, release models.Release) templateVars {
	return templateVars{
		RepoName:    repo.Repo,
		ReleaseBody: template.HTML(blackfriday.MarkdownBasic([]byte(release.Body))),
		ReleaseTag:  release.Tag,
		RTDomain:    os.Getenv("RT_DOMAIN"),
	}
}
