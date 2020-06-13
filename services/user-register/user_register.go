package useregister

import (
	"os"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	e "github.com/MihaiBlebea/purpletree/platform/email"
	u "github.com/MihaiBlebea/purpletree/platform/user"
)

// New creates a new ActivateTokenService
func New() *RegisterUserService {
	userRepository := *u.Repo(c.Mysql())
	emailService := e.Service{}

	return &RegisterUserService{userRepository, emailService}
}

// RegisterUserService is a service that registers a user in the platform
type RegisterUserService struct {
	userRepository u.Repository
	EmailService   e.Service
}

// RegisterUserResponse is the response struct for RegisterUserService
type RegisterUserResponse struct {
	UserID  int  `json:"user_id"`
	Success bool `json:"success"`
}

// Execute runs the RegisterUserService
func (s *RegisterUserService) Execute(name, email, password string, consent bool) (response RegisterUserResponse, err error) {
	user, err := u.New(name, email, password, consent)
	if err != nil {
		return response, err
	}

	// Mark user as not active until it confirms the email address
	user.Active = false

	// Generate confirm code for user
	user.ConfirmCode = u.GenerateRandomPassword()

	userID, err := s.userRepository.Add(user)
	if err != nil {
		return response, err
	}

	// After saving the user, send a confirmation email
	data := make(map[string]interface{})
	data["name"] = user.Name
	data["email"] = user.Email
	data["confirmUrl"] = os.Getenv("CLIENT_DO_ENDPOINT") + "?work=confirm-email&code=" + user.ConfirmCode

	s.EmailService.Send("confirm-email", data)

	return RegisterUserResponse{userID, true}, err
}
