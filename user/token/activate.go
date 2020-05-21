package token

import (
	c "github.com/MihaiBlebea/Wordpress/platform/connection"
)

// ActivateResponse struct
type ActivateResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
	Host    string `json:"host"`
	Active  bool   `json:"active"`
	Expires string `json:"expires"`
	Created string `json:"created"`
}

// Activate activates a token
func Activate(token, host string) (ActivateResponse, error) {
	tokenRepo := Repo(c.Mysql())
	tkn, _, err := tokenRepo.FindToken(token)
	if err != nil {
		return ActivateResponse{}, err
	}
	if tkn.HasHost() == false {
		tkn.Link(host)
		_, err = tokenRepo.Update(tkn)
		if err != nil {
			return ActivateResponse{}, err
		}
	}
	return NewActiveResponse("Token was activated", tkn), nil
}

// Deactivate deactivates a token
func Deactivate(token, host string) (ActivateResponse, error) {
	tokenRepo := Repo(c.Mysql())
	tkn, _, err := tokenRepo.FindToken(token)
	if err != nil {
		return ActivateResponse{}, err
	}
	if host == tkn.LinkedHost {
		tkn.Unlink()
		_, err = tokenRepo.Update(tkn)
		if err != nil {
			return ActivateResponse{}, err
		}
	}
	return NewActiveResponse("Token was deactivated", tkn), nil
}

// NewActiveResponse create new response
func NewActiveResponse(message string, token *Token) ActivateResponse {
	return ActivateResponse{
		message,
		token.Token,
		token.LinkedHost,
		token.Active,
		token.Expires,
		token.Created.Format(timeFormat),
	}
}
