package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const apiKey = "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"
const url = "https://api.sendinblue.com/v3/smtp/email"

// Sender struct
type Sender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Receiver struct
type Receiver struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// SendInBlueProvider is an email provider
type SendInBlueProvider struct {
	sender    Sender
	receivers []Receiver
	html      string
	subject   string
}

// SetSender set the sender property
func (p *SendInBlueProvider) SetSender(name, email string) {
	p.sender = Sender{name, email}
}

// SetReceiver set the receiver property
func (p *SendInBlueProvider) SetReceiver(name, email string) {
	p.receivers = append(p.receivers, Receiver{name, email})
}

// SetHTML set the html property
func (p *SendInBlueProvider) SetHTML(html string) {
	p.html = html
}

// SetSubject set the subject property
func (p *SendInBlueProvider) SetSubject(subject string) {
	p.subject = subject
}

// Send an email
func (p *SendInBlueProvider) Send() (string, error) {
	type Body struct {
		Sender      Sender     `json:"sender"`
		To          []Receiver `json:"to"`
		Subject     string     `json:"subject"`
		HTMLContent string     `json:"htmlContent"`
	}

	body := Body{
		Sender:      p.sender,
		To:          p.receivers,
		Subject:     p.subject,
		HTMLContent: p.html,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	request.Header.Set("Content-type", "application/json")
	request.Header.Set("api-key", apiKey)
	if err != nil {
		return "", err
	}

	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	resp, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(resp))
	return "", nil
}
