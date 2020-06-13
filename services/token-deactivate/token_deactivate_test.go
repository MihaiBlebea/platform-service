package tkndeactiv

import (
	"testing"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	tkn "github.com/MihaiBlebea/purpletree/platform/user/token"
)

func TestTokenIsDeactivated(t *testing.T) {
	var (
		email  = "john.doe@gmail.com"
		userID = 1
		host   = "google.com"
	)
	tokenRepository := *tkn.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})

	token, err := tkn.New(userID, email)
	if err != nil {
		t.Error(err)
	}
	token.LinkedHost = host

	id, err := tokenRepository.Add(token)
	if err != nil {
		t.Error(err)
	}
	token.ID = id

	service := DeactivateTokenService{tokenRepository}
	response, err := service.Execute(token.Token, host)
	if err != nil {
		t.Error(err)
	}

	if response.Token != token.Token {
		t.Errorf("Token output should match token input")
	}
	if response.Active != true {
		t.Errorf("Token should have been active")
	}
	if response.Host != "" {
		t.Errorf("Token should have host empty string, but got %s", response.Host)
	}

	t.Cleanup(func() {
		tokenRepository.RemoveByUser(userID)
	})
}
