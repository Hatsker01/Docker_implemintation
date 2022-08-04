package mail

import (
	"crypto/tls"
	gomail "gopkg.in/mail.v2"
)

func SendMail(num string, mail string) error {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "jamshidbek1805@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", mail)

	// Set E-Mail subject
	m.SetHeader("Subject", "Gomail test subject")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", num)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "jamshidbek1805@gmail.com", "qzsxlgudmntzsqhs")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	err := d.DialAndSend(m)

	return err
}
