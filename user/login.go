package user

import (
	"errors"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	"golang.org/x/crypto/bcrypt"
)

// LoginResponse struct
type LoginResponse struct {
	Message string `json:"message"`
	JWT     string `json:"jwt_token"`
	Success bool   `json:"success"`
}

// Login auth flow
func Login(email, password string) (LoginResponse, error) {
	userRepo := Repo(c.Mysql())
	user, count, err := userRepo.FindByEmail(email)
	if err != nil {
		return failResponse(), err
	}
	if count == 0 {
		return failResponse(), err
	}

	isPasswordValid := CheckPasswordHash(password, user.Password)
	if isPasswordValid == false {
		return failResponse(), err
	}

	return successResponse(user), nil
}

// Authenticate checks if a JWT is valid or not and returns bool
func Authenticate(token string) (bool, *User, error) {
	userRepo := Repo(c.Mysql())
	user, count, err := userRepo.FindByJWT(token)
	if err != nil {
		return false, &User{}, nil
	}
	if count == 0 {
		return false, &User{}, errors.New("No user found")
	}

	return true, user, nil
}

// CheckPasswordHash checks if the password is valid
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func successResponse(user *User) LoginResponse {
	return LoginResponse{
		JWT:     user.JWT,
		Message: "Authentication successfull",
		Success: true,
	}
}

func failResponse() LoginResponse {
	return LoginResponse{
		JWT:     "",
		Message: "Authentication failed",
		Success: false,
	}
}
