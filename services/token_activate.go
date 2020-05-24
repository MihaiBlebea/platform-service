package services

import tkn "github.com/MihaiBlebea/Wordpress/platform/user/token"

// ActivateTokenService activates a token, linking it with a host
type ActivateTokenService struct {
	TokenRepository tkn.Repository
}

// ActivateTokenResponse resonse for ActivateTokenService
type ActivateTokenResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Active  bool   `json:"active"`
	Expires string `json:"expires"`
	Created string `json:"created"`
}

// Execute runs the ActivateTokenService
func (s *ActivateTokenService) Execute(token, host string) (response ActivateTokenResponse, err error) {
	const timeFormat = "2006-01-02 15:04:05"

	tkn, _, err := s.TokenRepository.FindToken(token)
	if err != nil {
		return response, err
	}
	if tkn.HasHost() == false {
		tkn.Link(host)
		_, err = s.TokenRepository.Update(tkn)
		if err != nil {
			return response, err
		}
	}
	return ActivateTokenResponse{
		"Token was activated",
		true,
		tkn.Token,
		tkn.LinkedHost,
		tkn.Active,
		tkn.Expires,
		tkn.Created.Format(timeFormat),
	}, nil
}
