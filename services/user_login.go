package services

import (
	u "github.com/MihaiBlebea/Wordpress/platform/user"
)

// LoginUserService is a service that logs in a user in the platform
type LoginUserService struct {
	userRepository u.Repository
}

// LoginUserResponse is the response struct for LoginUserService
type LoginUserResponse struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
	JWT     string `json:"jwt_token"`
	Success bool   `json:"success"`
}

// Execute runs the RegisterUserService
func (s *LoginUserService) Execute(email, password string) (response LoginUserResponse, err error) {
	user, count, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return response, err
	}
	if count == 0 {
		return s.failResponse(), err
	}

	isPasswordValid := u.CheckPasswordHash(password, user.Password)
	if isPasswordValid == false {
		return s.failResponse(), err
	}

	return s.successResponse(user), nil
}

func (s *LoginUserService) successResponse(user *u.User) LoginUserResponse {
	return LoginUserResponse{
		ID:      user.ID,
		JWT:     user.JWT,
		Message: "Authentication successfull",
		Success: true,
	}
}

func (s *LoginUserService) failResponse() LoginUserResponse {
	return LoginUserResponse{
		Message: "Authentication failed",
		Success: false,
	}
}
