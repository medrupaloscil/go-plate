package services

import (
	"net/smtp"
	"os"
)

func SendEmail(address string, subject string, content string) (bool, error) {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")

	to := []string{
		address,
	}
	smtpHost := os.Getenv("EMAIL_HOST")
	smtpPort := os.Getenv("EMAIL_PORT")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	msg := []byte("From: " + from + "\r\n" +
		"To: " + address + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		content + "\r\n")
	
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		return false, err
	}

	return true, nil
}
