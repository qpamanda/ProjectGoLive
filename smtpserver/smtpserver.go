package smtpserver

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"

	"golang.org/x/crypto/bcrypt"
)

var (
	HostPath      string
	SMTPHost      string
	SMTPPort      string
	EmailPassword string
	FromEmail     string
)

// SendResetPwdEmail sends an email to user to reset their password
// Author: Amanda
func EmailResetPassword(email string) error {

	var body bytes.Buffer

	// Recipient email
	toEmail := []string{
		email,
	}

	// Host authentication
	auth := smtp.PlainAuth("", FromEmail, EmailPassword, SMTPHost)

	t, _ := template.ParseFiles("templates/emailresetpwd.html")

	// Hashed the email from user input
	bEmail, _ := bcrypt.GenerateFromPassword([]byte(email), bcrypt.MinCost)
	resetLink := HostPath + "/resetpwd?key=" + string(bEmail)

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Giving Grace Portal - Reset Password Request \n%s\n\n", mimeHeaders)))

	t.Execute(&body, struct {
		ResetLink string
	}{
		resetLink,
	})

	// Sending email.
	err := smtp.SendMail(SMTPHost+":"+SMTPPort, auth, FromEmail, toEmail, body.Bytes())
	if err != nil {
		return errors.New("failed to send reset password email")
	}
	return nil
}
