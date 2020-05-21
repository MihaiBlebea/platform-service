package token

import (
	c "github.com/MihaiBlebea/Wordpress/platform/connection"
)

// ValidateResponse response struct
type ValidateResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// Validate token
func Validate(token, host string) (ValidateResponse, error) {
	tokenRepo := Repo(c.Mysql())
	tkn, _, err := tokenRepo.FindToken(token)
	if err != nil {
		return ValidateResponse{}, err
	}

	valid, message, err := tkn.Validate(token, host)
	if err != nil {
		return ValidateResponse{}, err
	}

	if valid != true {
		return ValidateResponse{false, message}, err
	}

	return ValidateResponse{true, message}, nil
}
