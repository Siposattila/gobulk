package email

import (
	"net"
	"net/smtp"
	"regexp"
	"strings"

	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/logger"
)

func (e *Email) verifyEmail() uint8 {
	domain := e.Email[strings.LastIndex(e.Email, "@")+1:]

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		logger.LogError(err)

		return interfaces.EMAIL_INVALID
	}

	mxHost := mxRecords[0].Host

	client, err := smtp.Dial(mxHost + ":25")
	if err != nil {
		logger.LogError(err)

		return interfaces.EMAIL_INVALID
	}
	defer client.Close()

	client.Hello("gobulk.com")
	client.Mail("info@gobulk.com")
	rcptErr := client.Rcpt(e.Email)
	client.Quit()

	if rcptErr != nil {
		logger.LogError(rcptErr)

		return interfaces.EMAIL_INVALID
	}

	return interfaces.EMAIL_VALID
}

func IsEmail(email *string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, *email)

	return match
}
