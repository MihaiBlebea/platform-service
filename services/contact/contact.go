package contact

import (
	e "github.com/MihaiBlebea/purpletree/platform/email"
)

// New returns a new contact Service
func New() *Service {
	emailService := e.Service{}

	return &Service{emailService}
}

// Service send an email based on a contact request received
type Service struct {
	EmailService e.Service
}

// Response is a response for the Service
type Response struct {
	Success bool `json:"success"`
}

// Execute runs the contact Service
func (s *Service) Execute(name, email, subject, message string) (response Response, err error) {
	// Send the message
	data := make(map[string]interface{})
	data["email"] = "mihaiserban.blebea@gmail.com" // replace
	data["user_name"] = name
	data["user_email"] = email
	data["subject"] = subject
	data["message"] = message

	err = s.EmailService.Send("contact", data)
	if err != nil {
		return response, err
	}

	// Return response
	return Response{true}, nil
}
