package main

import (
	"crypto/tls"
	"errors"
	"log"
	"net"
	"net/smtp"
	"strings"
)

func main() {
	log.Println(verifyEmail("gobulk2024@outlook.com"))
}

func verifyEmail(email string) uint8 {
	domain := email[strings.LastIndex(email, "@")+1:]

	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Println(err)

		return 2
	}

	mxHost := mxRecords[0].Host
	mxHost = "smtp-mail.outlook.com"

	tlsconfig := &tls.Config{
		ServerName: mxHost,
	}

	client, err := smtp.Dial(mxHost + ":25")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	//conn, err := net.Dial("tcp", mxHost+":25")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer conn.Close()

	//client, err := smtp.NewClient(conn, mxHost)
	//if err != nil {
	//	log.Fatal(err)
	//}

	// StartTLS to upgrade the connection to TLS
	if err := client.StartTLS(tlsconfig); err != nil {
		log.Fatal(err)
	}

	client.Hello("gobulk.com")
	//auth := smtp.PlainAuth("", "gobulk2024@gmail.com", "xU!@'nWg85Ga-er", mxHost)
	auth := LoginAuth("gobulk2024@outlook.com", "xU!@'nWg85Ga-er")
	if err := client.Auth(auth); err != nil {
		log.Fatal(err)
	}

	client.Mail("info@gobulk.com")
	rcptErr := client.Rcpt(email)
	vrfyErr := client.Verify(email)
	client.Quit()

	if rcptErr != nil {
		log.Println(rcptErr)

		return 2
	}

	if vrfyErr != nil {
		log.Println(vrfyErr)

		return 2
	}

	return 1
}

type loginAuth struct {
	username, password string
}

func LoginAuth(username string, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown from server")
		}
	}
	return nil, nil
}
