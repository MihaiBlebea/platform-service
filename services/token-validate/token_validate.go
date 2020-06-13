package tknvalid

import (
	c "github.com/MihaiBlebea/purpletree/platform/connection"
	tkn "github.com/MihaiBlebea/purpletree/platform/user/token"
)

// New creates a new ActivateTokenService
func New() *ValidateTokenService {
	tokenRepository := *tkn.Repo(c.Mysql())
	return &ValidateTokenService{tokenRepository}
}

// ValidateTokenService validates a token by host
type ValidateTokenService struct {
	tokenRepository tkn.Repository
}

// ValidateTokenResponse response for ValidateTokenService
type ValidateTokenResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// Execute runs the token validation
func (s *ValidateTokenService) Execute(token, host string) (response ValidateTokenResponse) {
	tkn, _, err := s.tokenRepository.FindToken(token)
	if err != nil {
		return s.invalidResponse(err.Error())
	}

	valid, err := tkn.Validate(token, host)
	if err != nil {
		return s.invalidResponse(err.Error())
	}

	if valid != true {
		return s.invalidResponse(err.Error())
	}

	return s.validResponse()
}

func (s *ValidateTokenService) validResponse() (response ValidateTokenResponse) {
	return ValidateTokenResponse{true, "Token is valid"}

}

func (s *ValidateTokenService) invalidResponse(message string) (response ValidateTokenResponse) {
	return ValidateTokenResponse{false, message}
}
