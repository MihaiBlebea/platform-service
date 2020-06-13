package tknvalid

import (
	"testing"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	tkn "github.com/MihaiBlebea/purpletree/platform/user/token"
)

func TestTokenInvalidHost(t *testing.T) {
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

	service := ValidateTokenService{tokenRepository}
	response := service.Execute(token.Token, host)

	if response.Valid != false {
		t.Errorf("Expected token to be invalid")
	}

	t.Cleanup(func() {
		tokenRepository.RemoveByUser(userID)
	})
}

func TestTokenNotActive(t *testing.T) {
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
	token.Active = false
	token.LinkedHost = host

	id, err := tokenRepository.Add(token)
	if err != nil {
		t.Error(err)
	}
	token.ID = id

	service := ValidateTokenService{tokenRepository}
	response := service.Execute(token.Token, host)

	if response.Valid != false {
		t.Errorf("Expected token to be invalid")
	}

	t.Cleanup(func() {
		tokenRepository.RemoveByUser(userID)
	})
}

func TestTokenValidHost(t *testing.T) {
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

	service := ValidateTokenService{tokenRepository}
	response := service.Execute(token.Token, host)

	if response.Valid == false {
		t.Errorf("Expected be invalid")
	}

	t.Cleanup(func() {
		tokenRepository.RemoveByUser(userID)
	})
}
