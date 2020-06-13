package usrtkns

import (
	c "github.com/MihaiBlebea/purpletree/platform/connection"
	t "github.com/MihaiBlebea/purpletree/platform/user/token"
)

// New creates a new ActivateTokenService
func New() *GetUserTokensService {
	tokenRepository := *t.Repo(c.Mysql())
	return &GetUserTokensService{tokenRepository}
}

// GetUserTokensService gets all tokens for a user
type GetUserTokensService struct {
	TokenRepository t.Repository
}

// GetUserTokensResponse response for GetUserTokensService
type GetUserTokensResponse struct {
	UserID int     `json:"user_id"`
	Tokens []Token `json:"tokens"`
}

// Token represents an activation token
type Token struct {
	Token string `json:"token"`
	Host  string `json:"host"`
}

// Execute runs the RegisterUserService
func (s *GetUserTokensService) Execute(userID int) (response GetUserTokensResponse, err error) {
	tokens, count, err := s.TokenRepository.FindByUserID(userID)
	if err != nil {
		return response, err
	}

	if count == 0 {
		return GetUserTokensResponse{userID, []Token{}}, nil
	}

	var tkns []Token
	for _, token := range *tokens {
		tkns = append(tkns, Token{
			Token: token.Token,
			Host:  token.LinkedHost,
		})
	}

	return GetUserTokensResponse{userID, tkns}, nil
}
