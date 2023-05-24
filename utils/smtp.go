package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func SendEmail(to []string, from, subject string) error {

	// Get SMTP credentials from environment variables
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpUsername := os.Getenv("SMTP_USERNAME")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Convert smtpPort from string to integer
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		errMsg := fmt.Sprintf("failed to convert SMTP port to integer: %v", err)
		fmt.Println(errMsg)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetAddressHeader("Cc", "hamza.imtiaz47022@gmail.com", "Hamza Imtiaz")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", "Hello <b>Your feedback is pending</b> . Please click here to provide feedback for  <i>Employee 1</i>!")
	// m.Attach("/home/Alex/lolcat.jpg")

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return nil
}
