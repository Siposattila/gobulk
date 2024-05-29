package email

import (
	"bufio"
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/Siposattila/gobulk/internal/interfaces"
	"github.com/Siposattila/gobulk/internal/logger"
)

type body struct {
	subject     string
	greeting    string
	message     string
	farewell    string
	company     string
	unsubscribe string
}

func (b *body) GetSubject() string     { return b.subject }
func (b *body) GetGreeting() string    { return b.greeting }
func (b *body) GetMessage() string     { return b.message }
func (b *body) GetFarewell() string    { return b.farewell }
func (b *body) GetCompany() string     { return b.company }
func (b *body) GetUnsubscribe() string { return b.unsubscribe }

type client struct {
	auth  *smtp.Auth
	host  string
	port  string
	email string
	body  interfaces.EmailBodyInterface
}

func (c *client) SetEmailBody(body interfaces.EmailBodyInterface) { c.body = body }

func NewClient(dsn string, body interfaces.EmailBodyInterface) interfaces.EmailClientInterface {
	match, _ := regexp.MatchString(`^smtp:\/\/[^:]+:[^@]+@[^:]+:\d+$`, dsn)
	if !match {
		logger.Fatal("Bad email DSN!")
	}

	helper := dsn
	helper = helper[strings.LastIndex(helper, "smtp://")+7:]
	from := helper[:strings.Index(helper, ":")]
	password := helper[strings.Index(helper, ":")+1 : strings.LastIndex(helper, "@")]
	host := helper[strings.LastIndex(helper, "@")+1 : strings.LastIndex(helper, ":")]
	port := helper[strings.LastIndex(helper, ":")+1:]

	auth := smtp.PlainAuth("", from, password, host)

	return &client{
		auth:  &auth,
		host:  host,
		port:  port,
		email: from,
		body:  body,
	}
}

func NewBodyConsole(company string, unsubscribe string) interfaces.EmailBodyInterface {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please give a subject:")
	subject, err := reader.ReadString('\n')
	if err != nil {
		logger.Fatal(err)
	}
	subject = strings.TrimSpace(subject)

	fmt.Println("Please give a greeting:")
	greeting, err := reader.ReadString('\n')
	if err != nil {
		logger.Fatal(err)
	}
	greeting = strings.TrimSpace(greeting)

	fmt.Println("Please give a content:")
	message, err := reader.ReadString('\n')
	if err != nil {
		logger.Fatal(err)
	}
	message = strings.TrimSpace(message)

	fmt.Println("Please give a farewell:")
	farewell, err := reader.ReadString('\n')
	if err != nil {
		logger.Fatal(err)
	}
	farewell = strings.TrimSpace(farewell)

	return &body{
		subject:     subject,
		greeting:    greeting,
		message:     message,
		farewell:    farewell,
		company:     company,
		unsubscribe: unsubscribe,
	}
}

func NewBody(subject string, greeting string, message string, farewell string, company string, unsubscribe string) interfaces.EmailBodyInterface {
	return &body{
		subject:     subject,
		greeting:    greeting,
		message:     message,
		farewell:    farewell,
		company:     company,
		unsubscribe: unsubscribe,
	}
}

func (c *client) Send(e interfaces.EmailInterface) {
	template, error := template.ParseFiles("email.html")
	if error != nil {
		logger.Fatal("Failed to load template!")
	}

	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \nFrom: %s \nTo: %s \n%s\n\n", c.body.GetSubject(), c.body.GetCompany()+" <"+c.email+">", e.GetEmail(), mimeHeaders)))

	unsubscribeUrl := c.body.GetUnsubscribe() + "/" + e.GetEmail()
	greeting := strings.ReplaceAll(c.body.GetGreeting(), "[name]", e.GetName())
	template.Execute(&body, struct {
		Greeting    string
		Message     string
		Farewell    string
		Company     string
		Unsubscribe string
	}{
		Greeting:    greeting,
		Message:     c.body.GetMessage(),
		Farewell:    c.body.GetFarewell(),
		Company:     c.body.GetCompany(),
		Unsubscribe: unsubscribeUrl,
	})

	err := smtp.SendMail(c.host+":"+c.port, *c.auth, c.email, []string{e.GetEmail()}, body.Bytes())
	if err != nil {
		logger.Fatal("An error occured during email sending: " + err.Error())
	}
	e.SetSendStatus(interfaces.EMAIL_SEND_STATUS_SENT)
	logger.LogSuccess("Email sent to " + e.GetEmail() + ".")
}
