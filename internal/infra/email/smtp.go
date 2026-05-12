package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
)

func RenderEmailTemplate(templateName string, data any) (string, error) {
	exePath, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := filepath.Join(exePath, "internal/infra/email/templates", templateName)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var body bytes.Buffer

	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}

	return body.String(), nil
}

func SendEmail(to, subject, body string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	from := os.Getenv("SMTP_FROM")

	msg := []byte(
		"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	return smtp.SendMail(
		host+":"+port,
		nil,
		from,
		[]string{to},
		msg,
	)
}
