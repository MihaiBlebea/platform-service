package tknactiv

import (
	"testing"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	tkn "github.com/MihaiBlebea/Wordpress/platform/user/token"
)

func TestTokenIsActivated(t *testing.T) {
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

	id, err := tokenRepository.Add(token)
	if err != nil {
		t.Error(err)
	}
	token.ID = id

	service := ActivateTokenService{tokenRepository}
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
	if response.Host != host {
		t.Errorf("Token should have host %s, but got %s", host, response.Host)
	}

	t.Cleanup(func() {
		tokenRepository.RemoveByUser(userID)
	})
}
