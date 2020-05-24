package usregister

import (
	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	u "github.com/MihaiBlebea/Wordpress/platform/user"
)

// New creates a new ActivateTokenService
func New() *RegisterUserService {
	userRepository := *u.Repo(c.Mysql())
	return &RegisterUserService{userRepository}
}

// RegisterUserService is a service that registers a user in the platform
type RegisterUserService struct {
	userRepository u.Repository
}

// RegisterUserResponse is the response struct for RegisterUserService
type RegisterUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	JWT  string `json:"jwt_token"`
}

// Execute runs the RegisterUserService
func (s *RegisterUserService) Execute(name, email, password string, consent bool) (response RegisterUserResponse, err error) {
	user, err := u.New(name, email, password, consent)
	if err != nil {
		return response, err
	}

	userID, err := s.userRepository.Add(user)
	if err != nil {
		return response, err
	}

	return RegisterUserResponse{userID, user.Name, user.JWT}, err
}
