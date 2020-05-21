package user

import (
	c "github.com/MihaiBlebea/Wordpress/platform/connection"
)

// RegisterResponse struct
type RegisterResponse struct {
	Name string `json:"name"`
	JWT  string `json:"jwt_token"`
}

// Register generates a user entry and JWT tokens
func Register(name, email, password string, consent bool) (RegisterResponse, error) {
	userRepo := Repo(c.Mysql())

	user, err := New(name, email, password, consent)
	if err != nil {
		return RegisterResponse{}, err
	}

	userID, err := userRepo.Add(user)
	if err != nil {
		return RegisterResponse{}, err
	}
	user.ID = userID

	return NewRegisterResponse(user), nil
}

// NewRegisterResponse creates an response object
func NewRegisterResponse(user *User) RegisterResponse {
	response := RegisterResponse{user.Name, user.JWT}
	return response
}
