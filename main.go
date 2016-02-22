package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"os/exec"
	"strings"
	"time"
)

const dateFormat = "02.01.2006. 15:04:05"

func main() {

	toEmail := os.Getenv("SU_TO_EMAIL")

	fromName := os.Getenv("SU_FROM_NAME")
	fromEmail := os.Getenv("SU_FROM_EMAIL")

	username := os.Getenv("SU_USER")
	pass := os.Getenv("SU_PASS")
	host := os.Getenv("SU_HOST")
	port := os.Getenv("SU_PORT")

	verbose := strings.ToLower(os.Getenv("SU_VERBOSE"))

	var result []byte
	defer func(start time.Time) {

		end := time.Now()
		elapsed := time.Since(start)

		result = []byte(fmt.Sprintf(
			"Start: %s\nEnd: %s\nElapsed: %s\n\n%s",
			start.Format(dateFormat),
			end.Format(dateFormat),
			elapsed,
			result))

		sendMail(result, toEmail, fromName, fromEmail, username, pass, host, port)

	}(time.Now())

	var params []string
	switch verbose {
	case "true", "on", "verbose", "1", "enable", "enabled":
		params = []string{"-v", "--apt-debug"}
	}

	cmd := exec.Command("unattended-upgrade", params...)
	result, err := cmd.CombinedOutput()
	if err != nil {
		result = []byte(fmt.Sprintf("ERROR: %s\n\n%s", err, result))
	}

}

func sendMail(result []byte, toEmail, fromName, fromEmail, username, pass, host, port string) {

	from := mail.Address{Name: fromName, Address: fromEmail}
	to := mail.Address{Name: toEmail, Address: toEmail}

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("os.Hostname failed with %q", err)
	}

	subject := fmt.Sprintf("[%s] unattended-upgrade report", hostname)
	body := string(result)

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
	defer func() {
		err = c.Quit()
		if err != nil {
			log.Fatalf("smtp.Client.Quit failed with %q", err)
		}
	}()

	// Auth
	if err = c.Auth(auth); err != nil {
		log.Fatalf("smtp.Client.Auth failed with %q", err)
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		log.Fatalf("smtp.Client.Mail failed with %q", err)
	}

	if err = c.Rcpt(to.Address); err != nil {
		log.Fatalf("smtp.Client.Rcpt failed with %q", err)
	}

	// Data
	w, err := c.Data()
	if err != nil {
		log.Fatalf("smtp.Client.Data failed with %q", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Fatalf("Msg Writer Write failed with %q", err)
	}

	err = w.Close()
	if err != nil {
		log.Fatalf("Msg Writer Close failed with %q", err)
	}
}
