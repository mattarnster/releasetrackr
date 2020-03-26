package emails

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"releasetrackr/models"

	"github.com/russross/blackfriday"

	"github.com/mailjet/mailjet-apiv3-go"
)

type templateVars struct {
	RepoName    string
	ReleaseBody template.HTML
	ReleaseTag  string
	ReleaseURL  string
	RTDomain    string
}

// SendNotificationEmail does exactly what it says on the tin.
func SendNotificationEmail(repo models.Repo, email string, release models.Release) {
	t, _ := template.ParseFiles("templates/new-release.html")
	var doc bytes.Buffer
	err := t.Execute(&doc, generateTemplateVars(repo, release))

	if err != nil {
		log.Printf("Template parse failed: %v", err.Error())
	}

	mj := mailjet.NewMailjetClient(os.Getenv("MAILJET_API_PUBLIC_KEY"), os.Getenv("MAILJET_API_PRIVATE_KEY"))
	messageWrapper := []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: "releasetrackr@mattarnster.co.uk",
				Name:  "releasetrackr",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: email,
				},
			},
			Subject:  "New release for " + repo.Repo,
			TextPart: "Hey there, A new release for " + repo.Repo + " has been detected!",
			HTMLPart: doc.String(),
		},
	}

	messages := mailjet.MessagesV31{Info: messageWrapper}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[Helper][NotificationEmail] Response: %+v\n", res)
}

// This function returns a struct of the required template variables.
func generateTemplateVars(repo models.Repo, release models.Release) templateVars {
	return templateVars{
		RepoName:    repo.Repo,
		ReleaseBody: template.HTML(blackfriday.MarkdownBasic([]byte(release.Body))),
		ReleaseTag:  release.Tag,
		ReleaseURL:  release.URL,
		RTDomain:    os.Getenv("RT_DOMAIN"),
	}
}
