package helpers

import (
	"bytes"
	"html/template"
	"log"
	"os"

	mailgun "gopkg.in/mailgun/mailgun-go.v1"
)

type verificationVars struct {
	VerificationLink string
}

// SendVerificationEmail does exactly what it says on the tin.
func SendVerificationEmail(email string, vt string) {
	mg := mailgun.NewMailgun("mattarnster.co.uk", os.Getenv("MAILGUN_API_KEY"), "")
	message := mailgun.NewMessage(
		"releasetrackr <releasetrackr@mattarnster.co.uk>",
		"Verify your email",
		"Hey there, please visit "+os.Getenv("RT_DOMAIN")+"/verify?key="+vt+" to verify your email address so that you can receive releasetrackr notifications!",
		email)

	t, _ := template.ParseFiles("templates/user-verification.html")
	var doc bytes.Buffer
	err := t.Execute(&doc, generateTemplateVarsVerification(vt))

	if err != nil {
		log.Printf("Template parse failed: %v", err.Error())
	}

	message.SetHtml(doc.String())

	resp, id, err := mg.Send(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[Helper][Mailgun] Queued ID: %s Resp: %s\n", id, resp)
}

// This function returns a struct filled with the required variables for use
// in the template
func generateTemplateVarsVerification(vt string) verificationVars {
	verificationLink := os.Getenv("RT_DOMAIN") + "/verify?key=" + vt
	return verificationVars{
		VerificationLink: verificationLink,
	}
}
