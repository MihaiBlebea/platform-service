package usrpassconfirm

import (
	"fmt"
	"os"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	e "github.com/MihaiBlebea/purpletree/platform/email"
	u "github.com/MihaiBlebea/purpletree/platform/user"
)

// New creates a new ResetUserPasswordConfirmService
func New() *ResetPasswordConfirmService {
	userRepository := *u.Repo(c.Mysql())
	emailService := e.Service{}
	return &ResetPasswordConfirmService{userRepository, emailService}
}

// ResetPasswordConfirmService resets the user password
type ResetPasswordConfirmService struct {
	UserRepository u.Repository
	EmailService   e.Service
}

// ResetPasswordConfirmResponse response for ResetUserPasswordService
type ResetPasswordConfirmResponse struct {
	Success bool `json:"success"`
}

// Execute runs the ResetPasswordConfirmService
func (s *ResetPasswordConfirmService) Execute(confirmEndpoint, email string) (response ResetPasswordConfirmResponse, err error) {
	// Try and find a user by email
	user, count, err := s.UserRepository.FindByEmail(email)
	if err != nil {
		return response, err
	}
	if count == 0 {
		return response, fmt.Errorf("No user found with email %s", email)
	}

	// Generate random confirmation code
	user.ConfirmCode = u.GenerateRandomPassword()
	_, err = s.UserRepository.Update(user)
	if err != nil {
		return response, err
	}

	// If user is found, then send him an email with a confirmation link
	data := make(map[string]interface{})
	data["name"] = user.Name
	data["email"] = user.Email
	data["confirmUrl"] = os.Getenv("CLIENT_DO_ENDPOINT") + "?work=confirm-password&code=" + user.ConfirmCode

	s.EmailService.Send("confirm-password", data)

	// Return response
	return ResetPasswordConfirmResponse{true}, nil
}
