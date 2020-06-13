package usremailvalid

import (
	"fmt"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	u "github.com/MihaiBlebea/purpletree/platform/user"
)

// New creates a new ResetUserPasswordConfirmService
func New() *ValidateUserEmailService {
	userRepository := *u.Repo(c.Mysql())

	return &ValidateUserEmailService{userRepository}
}

// ValidateUserEmailService validates if the user email is unique in the db
type ValidateUserEmailService struct {
	UserRepository u.Repository
}

// ValidateUserEmailResponse response for ValidateUserEmailService
type ValidateUserEmailResponse struct {
	Success bool `json:"success"`
}

// Execute runs the ResetPasswordConfirmService
func (s *ValidateUserEmailService) Execute(email string) (response ValidateUserEmailResponse, err error) {
	// Try and find a user by email
	_, count, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return response, err
	}
	if count > 0 {
		return response, fmt.Errorf("User already found with email %s", email)
	}

	// Return response
	return ValidateUserEmailResponse{true}, nil
}
