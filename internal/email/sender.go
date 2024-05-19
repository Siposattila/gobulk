package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"regexp"
	"strings"
	"text/template"

	"github.com/Siposattila/gobulk/internal/console"
)

type ClientInterface interface {
	Send(e *Email)
}

type client struct {
	Auth        *smtp.Auth
	Host        *string
	Port        *string
	Email       *string
	Subject     *string
	Message     *string
	Company     *string
	Unsubscribe *string
}

func NewClient(dsn *string, subject string, message string, company *string, unsubscribe *string) ClientInterface {
	match, _ := regexp.MatchString(`^smtp:\/\/[^:]+:[^@]+@[^:]+:\d+$`, *dsn)
	if !match {
		console.Fatal("Bad email DSN!")
	}

	helper := *dsn
	helper = helper[strings.LastIndex(helper, "smtp://")+7:]
	from := helper[:strings.Index(helper, ":")]
	password := helper[strings.Index(helper, ":")+1 : strings.LastIndex(helper, "@")]
	host := helper[strings.LastIndex(helper, "@")+1 : strings.LastIndex(helper, ":")]
	port := helper[strings.LastIndex(helper, ":")+1:]

	auth := smtp.PlainAuth("", from, password, host)

	return &client{
		Auth:        &auth,
		Host:        &host,
		Port:        &port,
		Email:       &from,
		Subject:     &subject,
		Message:     &message,
		Company:     company,
		Unsubscribe: unsubscribe,
	}
}

func (c *client) Send(e *Email) {
	template, error := template.ParseFiles("bulk.html")
	if error != nil {
		console.Fatal("Failed to load template!")
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \nFrom: %s \nTo: %s \n%s\n\n", *c.Subject, *c.Email, e.Email, mimeHeaders)))

	template.Execute(&body, struct {
		Name        string
		Message     string
		Company     string
		Unsubscribe string
	}{
		Name:        e.Name,
		Message:     *c.Message,
		Company:     *c.Company,
		Unsubscribe: *c.Unsubscribe,
	})

	err := smtp.SendMail(*c.Host+":"+*c.Port, *c.Auth, *c.Email, []string{e.Email}, body.Bytes())
	if err != nil {
		console.Fatal("An error occured during email sending: " + err.Error())
	}
	console.Success("Email sent to " + e.Email)
}
