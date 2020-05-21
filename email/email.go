package email

import (
	"bytes"
	"fmt"
	"text/template"
)

const senderName = "Mihai"
const senderEmail = "mihai@getchip.uk"

// Provider interface for email providers
type Provider interface {
	Send() (string, error)
	SetSender(name, email string)
	SetReceiver(name, email string)
	SetHTML(html string)
	SetSubject(subject string)
}

// Details for email sending
type Details struct {
	TemplateName  string
	Data          interface{}
	Subject       string
	ReceiverName  string
	ReceiverEmail string
}

// Send email
func Send(provider Provider, details Details) (string, error) {
	tmpl, err := template.ParseFiles(fmt.Sprintf("./email/templates/%s.html", details.TemplateName))
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err := tmpl.Execute(&tpl, details.Data); err != nil {
		return "", err
	}

	result := tpl.String()

	provider.SetSender(senderName, senderEmail)
	provider.SetReceiver(details.ReceiverName, details.ReceiverEmail)
	provider.SetHTML(result)
	provider.SetSubject(details.Subject)
	response, err := provider.Send()
	if err != nil {
		return "", err
	}
	return response, err
}
