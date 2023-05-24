package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Shopify/gomail"
)

// SendEmail sends an email using a remote SMTP server.
func SendEmail(toAddresses []string, fromAddress, subject string) error {
	// SMTP server credentials
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

	// Create a new email message
	message := gomail.NewMessage()
	message.SetHeader("From", fromAddress)
	message.SetHeader("To", toAddresses...)
	message.SetHeader("Subject", subject)

	// Set the HTML body
	message.SetBody("text/html", getHTMLBody())

	// Create the SMTP dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Send the email
	if err := dialer.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// getHTMLBody returns the HTML body content of the email
func getHTMLBody() string {
	// You can customize the HTML body template here
	htmlBody := `
		<!DOCTYPE html>
		<html>
			<head>
			 <meta charset="UTF-8">
             <meta http-equiv="X-UA-Compatible" content="IE=edge">
             <meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Email Template</title>
			</head>
			<body>
				<h1>Your feedback is pending</h1>
				<p>Please click <a href="https://asdemo.teo-intl.com/dashboard">here</a> to provide feedback for Employee 1.</p>
			</body>
		</html>
	`

	return htmlBody
}
