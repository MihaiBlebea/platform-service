package useremailconfirm

import (
	"fmt"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	e "github.com/MihaiBlebea/purpletree/platform/email"
	u "github.com/MihaiBlebea/purpletree/platform/user"
)

// New creates a new ResetUserPasswordConfirmService
func New() *ConfirmUserEmailService {
	userRepository := *u.Repo(c.Mysql())
	emailService := e.Service{}
	return &ConfirmUserEmailService{userRepository, emailService}
}

// ConfirmUserEmailService resets the user password
type ConfirmUserEmailService struct {
	UserRepository u.Repository
	EmailService   e.Service
}

// ConfirmUserEmailResponse response for ConfirmUserEmailService
type ConfirmUserEmailResponse struct {
	Success bool   `json:"success"`
	Name    string `json:"name"`
	JWT     string `json:"jwt_token"`
}

// Execute runs the ResetPasswordConfirmService
func (s *ConfirmUserEmailService) Execute(code string) (response ConfirmUserEmailResponse, err error) {
	// Try and find a user by code
	user, count, err := s.UserRepository.FindByConfirmCode(code)
	if err != nil {
		return response, err
	}
	if count == 0 {
		return response, fmt.Errorf("No user found with code %s", code)
	}

	user.Active = true
	user.ConfirmCode = ""

	_, err = s.UserRepository.Update(user)
	if err != nil {
		return response, err
	}

	// If user is found, then send him an email with a confirmation link
	data := make(map[string]interface{})
	data["name"] = user.Name
	data["email"] = user.Email

	s.EmailService.Send("welcome", data)

	// Return response
	return ConfirmUserEmailResponse{true, user.Name, user.JWT}, nil
}
