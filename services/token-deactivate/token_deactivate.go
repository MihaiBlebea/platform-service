package tkndeactiv

import (
	c "github.com/MihaiBlebea/purpletree/platform/connection"
	tkn "github.com/MihaiBlebea/purpletree/platform/user/token"
)

// New creates a new ActivateTokenService
func New() *DeactivateTokenService {
	tokenRepository := *tkn.Repo(c.Mysql())
	return &DeactivateTokenService{tokenRepository}
}

// DeactivateTokenService deactivates a token, removing the linked host
type DeactivateTokenService struct {
	TokenRepository tkn.Repository
}

// DeactivateTokenResponse resonse for DeactivateTokenService
type DeactivateTokenResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Active  bool   `json:"active"`
	Expires string `json:"expires"`
	Created string `json:"created"`
}

// Execute runs the ActivateTokenService
func (s *DeactivateTokenService) Execute(token, host string) (response DeactivateTokenResponse, err error) {
	const timeFormat = "2006-01-02 15:04:05"

	tkn, _, err := s.TokenRepository.FindToken(token)
	if err != nil {
		return response, err
	}
	if host == tkn.LinkedHost {
		tkn.Unlink()
		_, err = s.TokenRepository.Update(tkn)
		if err != nil {
			return response, err
		}
	}

	return DeactivateTokenResponse{
		"Token was deactivated",
		true,
		tkn.Token,
		tkn.LinkedHost,
		tkn.Active,
		tkn.Expires,
		tkn.Created.Format(timeFormat),
	}, nil
}
