package email

import (
	"net"
	"net/smtp"
	"regexp"
	"strings"

	"github.com/Siposattila/gobulk/internal/console"
)

func (e *Email) verifyEmail() bool {
	domain := e.Email[strings.LastIndex(e.Email, "@")+1:]

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		console.Error("No MX record found")

		return EMAIL_INVALID
	}

	mxHost := mxRecords[0].Host

	client, err := smtp.Dial(mxHost + ":25")
	if err != nil {
		console.Error("Failed to connect to SMTP server")

		return EMAIL_INVALID
	}
	defer client.Close()

	client.Hello("localhost")     // TODO: needs real domain
	client.Mail("me@example.com") // TODO: needs real email
	rcptErr := client.Rcpt(e.Email)
	client.Quit()

	if rcptErr != nil {
		console.Warning("Invalid email address")

		return EMAIL_INVALID
	}

	return EMAIL_VALID
}

func IsEmail(email *string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, *email)

	return match
}
