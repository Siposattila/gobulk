package email

import (
	"net"
	"net/smtp"
	"regexp"
	"strings"
)

func (e *Email) verifyEmail() bool {
	domain := e.Email[strings.LastIndex(e.Email, "@")+1:]

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		return EMAIL_INVALID
	}

	mxHost := mxRecords[0].Host

	client, err := smtp.Dial(mxHost + ":25")
	if err != nil {
		return EMAIL_INVALID
	}
	defer client.Close()

	client.Hello("gobulk.com")
	client.Mail("info@gobulk.com")
	rcptErr := client.Rcpt(e.Email)
	client.Quit()

	if rcptErr != nil {
		return EMAIL_INVALID
	}

	return EMAIL_VALID
}

func IsEmail(email *string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, *email)

	return match
}
