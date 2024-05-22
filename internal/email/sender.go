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

type EmailBody struct {
	Subject     string
	Greeting    string
	Message     string
	Farewell    string
	Company     string
	Unsubscribe string
}

type client struct {
	Auth  *smtp.Auth
	Host  string
	Port  string
	Email string
	Body  *EmailBody
}

func NewClient(dsn string, body *EmailBody) ClientInterface {
	match, _ := regexp.MatchString(`^smtp:\/\/[^:]+:[^@]+@[^:]+:\d+$`, dsn)
	if !match {
		console.Fatal("Bad email DSN!")
	}

	helper := dsn
	helper = helper[strings.LastIndex(helper, "smtp://")+7:]
	from := helper[:strings.Index(helper, ":")]
	password := helper[strings.Index(helper, ":")+1 : strings.LastIndex(helper, "@")]
	host := helper[strings.LastIndex(helper, "@")+1 : strings.LastIndex(helper, ":")]
	port := helper[strings.LastIndex(helper, ":")+1:]

	auth := smtp.PlainAuth("", from, password, host)

	return &client{
		Auth:  &auth,
		Host:  host,
		Port:  port,
		Email: from,
		Body:  body,
	}
}

func (c *client) Send(e *Email) {
	template, error := template.ParseFiles("bulk.html")
	if error != nil {
		console.Fatal("Failed to load template!")
	}

	emailBody := c.Body
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \nFrom: %s \nTo: %s \n%s\n\n", c.Body.Subject, c.Email, e.Email, mimeHeaders)))

	emailBody.Unsubscribe = c.Body.Unsubscribe + "/" + e.Email
	emailBody.Greeting = strings.ReplaceAll(c.Body.Greeting, "[name]", e.Name)
	template.Execute(&body, c.Body)

	err := smtp.SendMail(c.Host+":"+c.Port, *c.Auth, c.Email, []string{e.Email}, body.Bytes())
	if err != nil {
		console.Fatal("An error occured during email sending: " + err.Error())
	}
	e.SendStatus = EMAIL_SEND_STATUS_SENT
}
