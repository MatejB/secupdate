package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
)

func main() {

	toEmail := os.Getenv("SU_TO_EMAIL")

	fromName := os.Getenv("SU_FROM_NAME")
	fromEmail := os.Getenv("SU_FROM_EMAIL")

	username := os.Getenv("SU_USER")
	pass := os.Getenv("SU_PASS")
	host := os.Getenv("SU_HOST")
	port := os.Getenv("SU_PORT")

	from := mail.Address{Name: fromName, Address: fromEmail}
	to := mail.Address{Name: toEmail, Address: toEmail}

	hostname, err := os.Hostname()
	if err != nil {
		log.Panic(err)
	}

	subject := fmt.Sprintf("[%s] unattended-upgrade report", hostname)
	body := "This is an example body.\n With two lines."

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := fmt.Sprintf("%s:%s", host, port)

	auth := smtp.PlainAuth("", username, pass, host)

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		log.Fatalf("tls.Dial %q failed with %q", servername, err)
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Fatalf("smtp.NewClient failed with %q", err)
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Fatalf("Auth failed with %q", err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Fatalf("Mail failed with %q", err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Fatalf("Rcpt failed with %q", err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Fatalf("Data failed with %q", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Fatalf("Write failed with %q", err)
	}

	err = w.Close()
	if err != nil {
		log.Fatalf("Close failed with %q", err)
	}

	c.Quit()

}
