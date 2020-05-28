package email

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

// Service connects to the email-service microservice
type Service struct {
}

// Send sends an email to the microservice
func (s *Service) Send(kind string, data map[string]interface{}) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if os.Getenv("EMAIL_HOST") == "" || os.Getenv("EMAIL_PORT") == "" {
		return errors.New("Undefined env variable")
	}

	if _, ok := data["email"]; !ok {
		return errors.New("Data should contain key email")
	}

	baseURL := fmt.Sprintf("http://%s:%s/send/%s", os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT"), kind)

	response, err := http.Post(baseURL, "application/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	fmt.Println("Email sender response", response)

	return nil
}
