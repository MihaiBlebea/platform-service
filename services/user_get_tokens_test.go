package services

import (
	"testing"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	u "github.com/MihaiBlebea/Wordpress/platform/user"
	tkn "github.com/MihaiBlebea/Wordpress/platform/user/token"
)

func TestUserHasGotTokens(t *testing.T) {
	var (
		name     = "Serban"
		email    = "serban@gmail.com"
		password = "intrex"
		consent  = true
	)

	userRepository := *u.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})

	registerService := RegisterUserService{
		userRepository,
	}
	registerResponse, err := registerService.Execute(name, email, password, consent)
	if err != nil {
		t.Error(err)
	}

	tokenRepository := *tkn.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})

	tokens, err := tkn.NewCount(registerResponse.ID, email, 3)
	if err != nil {
		t.Error(err)
	}
	for _, token := range *tokens {
		tokenRepository.Add(&token)
	}

	service := GetUserTokensService{
		tokenRepository,
	}

	response, err := service.Execute(registerResponse.ID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Tokens) != 3 {
		t.Errorf("Expected to have %d tokens, but got %d", 3, len(response.Tokens))
	}
	if response.UserID != registerResponse.ID {
		t.Errorf("Expected user id to be %d, but got %d", registerResponse.ID, response.UserID)
	}

	t.Cleanup(func() {
		user := u.User{ID: response.UserID}
		userRepository.Remove(&user)

		tokenRepository.RemoveByUser(response.UserID)
	})
}

func TestUserHasNoTokens(t *testing.T) {
	var (
		name     = "Serban"
		email    = "serban@gmail.com"
		password = "intrex"
		consent  = true
	)

	userRepository := *u.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})

	registerService := RegisterUserService{
		userRepository,
	}
	registerResponse, err := registerService.Execute(name, email, password, consent)
	if err != nil {
		t.Error(err)
	}

	tokenRepository := *tkn.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})

	service := GetUserTokensService{
		tokenRepository,
	}

	response, err := service.Execute(registerResponse.ID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Tokens) != 0 {
		t.Errorf("Expected to have %d tokens, but got %d", 0, len(response.Tokens))
	}
	if response.UserID != registerResponse.ID {
		t.Errorf("Expected user id to be %d, but got %d", registerResponse.ID, response.UserID)
	}

	t.Cleanup(func() {
		user := u.User{ID: response.UserID}
		userRepository.Remove(&user)

		tokenRepository.RemoveByUser(response.UserID)
	})
}
