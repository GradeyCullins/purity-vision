package mail

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	Name    string
	To      string
	Subject string
	Plain   string
	Html    string
}

func SendMail(email Email) error {
	from := mail.NewEmail("GRAYWAREZ", "info@graywarez.com")

	to := mail.NewEmail(email.Name, email.To)
	message := mail.NewSingleEmail(from, email.Subject, to, email.Plain, email.Html)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	_, err := client.Send(message)
	if err != nil {
		return err
	} else {
		// log.Println(response.StatusCode)
		// log.Println(response.Body)
		// log.Println(response.Headers)
		return nil
	}

}

func SendLicenseMail(emailTo string, licenseID string) error {
	email := Email{
		Name:    emailTo,
		To:      emailTo,
		Subject: "Your Purity Vision License is Here!",
		Plain:   fmt.Sprintf("Your PurityVision License Key: %s\n", licenseID),
		Html:    fmt.Sprintf("<h1>Your PurityVision License Key</h1><p>%s</p>", licenseID),
	}

	if err := SendMail(email); err != nil {
		return err
	}

	return nil
}
