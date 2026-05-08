package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/resend/resend-go/v2"
)

type ResendMailer struct {
	fromEmail string
	client    *resend.Client
}

func NewResend(apiKey, fromEmail string) *ResendMailer {
	client := resend.NewClient(apiKey)

	return &ResendMailer{
		fromEmail: fromEmail,
		client:    client,
	}
}

func (m *ResendMailer) Send(templateFile, username, email string, data any, isSandbox bool) error {
	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	// Build the send email request
	request := &resend.SendEmailRequest{
		From:    m.fromEmail,
		To:      []string{email},
		Subject: subject.String(),
		Html:    body.String(),
	}

	// Retry logic with exponential backoff
	for i := 0; i < maxRetires; i++ {
		sent, err := m.client.Emails.Send(request)
		if err != nil {
			log.Printf("Failed to send email to %v, attempt %d of %d", email, i+1, maxRetires)
			log.Printf("Error: %v", err)

			// exponential backoff before retrying
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent successfully with ID: %v", sent.Id)
		return nil
	}

	return fmt.Errorf("failed to send email after %d attempts", maxRetires)
}
