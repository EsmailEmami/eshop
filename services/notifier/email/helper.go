package email

import (
	"bytes"
	"fmt"
	"text/template"
)

type emailHeader struct {
	From    string
	To      string
	Subject string
}

func (eh emailHeader) String() string {
	var buf bytes.Buffer

	buf.WriteString("From: " + eh.From + "\r\n")
	buf.WriteString("To: " + eh.To + "\r\n")
	buf.WriteString("Subject: " + eh.Subject + "\r\n")
	buf.WriteString("MIME-Version: " + "1.0" + "\r\n")
	buf.WriteString("Content-Type: " + "text/html; charset=\"utf-8\"" + "\r\n")

	return buf.String()
}

func generateMessage(key Key, data any, header *emailHeader) ([]byte, error) {
	templatePath := key.getTemplatePath()

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return nil, fmt.Errorf("Invalid template path. Error: %s", err.Error())
	}

	var emailBody bytes.Buffer

	emailBody.WriteString(header.String())

	if err := tmpl.Execute(&emailBody, data); err != nil {
		return nil, fmt.Errorf("Error while rendering the template. Error: %s", err.Error())
	}

	return emailBody.Bytes(), nil
}
