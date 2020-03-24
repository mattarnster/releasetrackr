package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"

	"github.com/mailjet/mailjet-apiv3-go"
)

type verificationVars struct {
	VerificationLink string
}

// SendVerificationEmail does exactly what it says on the tin.
func SendVerificationEmail(email string, vt string) {
	t, _ := template.ParseFiles("templates/user-verification.html")
	var doc bytes.Buffer
	err := t.Execute(&doc, generateTemplateVarsVerification(vt))

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
			Subject:  "Verify your email",
			TextPart: "Hey there, please visit " + os.Getenv("RT_DOMAIN") + "/verify?key=" + vt + " to verify your email address so that you can receive releasetrackr notifications!",
			HTMLPart: doc.String(),
		},
	}

	messages := mailjet.MessagesV31{Info: messageWrapper}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("[Helper][Mailjet] Response: %+v\n", res)
}

// This function returns a struct filled with the required variables for use
// in the template
func generateTemplateVarsVerification(vt string) verificationVars {
	verificationLink := os.Getenv("RT_DOMAIN") + "/verify?key=" + vt
	return verificationVars{
		VerificationLink: verificationLink,
	}
}
