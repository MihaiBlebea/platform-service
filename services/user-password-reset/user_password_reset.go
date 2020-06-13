package usrpassreset

import (
	"fmt"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	u "github.com/MihaiBlebea/purpletree/platform/user"
)

// New creates a new ResetUserPasswordService
func New() *ResetUserPasswordService {
	userRepository := *u.Repo(c.Mysql())
	return &ResetUserPasswordService{userRepository}
}

// ResetUserPasswordService resets the user password
type ResetUserPasswordService struct {
	UserRepository u.Repository
}

// ResetUserPasswordResponse response for ResetUserPasswordService
type ResetUserPasswordResponse struct {
	UserID   int    `json:"user_id"`
	JWTToken string `json:"jwt_token"`
	Success  bool   `json:"success"`
}

// Execute runs the ResetUserPasswordService
func (s *ResetUserPasswordService) Execute(confirmCode, newPassword string) (response ResetUserPasswordResponse, err error) {
	// Get user by email
	user, count, err := s.UserRepository.FindByConfirmCode(confirmCode)
	if err != nil {
		return response, err
	}
	if count == 0 {
		return response, fmt.Errorf("No user found with confirm code %s", confirmCode)
	}

	// Update the user password
	user.HashPassword(newPassword)

	// Remove confirm code
	user.ConfirmCode = ""

	// Update user JWT token
	err = user.GenerateJWT()
	if err != nil {
		return response, err
	}

	// Save updated user
	s.UserRepository.Update(user)

	// Return new user details
	return ResetUserPasswordResponse{
		user.ID,
		user.JWT,
		true,
	}, nil
}
