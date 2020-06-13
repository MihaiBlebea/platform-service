package usrtkns

import (
	"testing"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	u "github.com/MihaiBlebea/purpletree/platform/user"
	tkn "github.com/MihaiBlebea/purpletree/platform/user/token"
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
	user := &u.User{
		Name:     name,
		Email:    email,
		Password: password,
		Consent:  consent,
	}
	id, err := userRepository.Add(user)
	if err != nil {
		t.Error(err)
	}
	user.ID = id

	tokenRepository := *tkn.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})

	tokens, err := tkn.NewCount(user.ID, email, 3)
	if err != nil {
		t.Error(err)
	}
	for _, token := range *tokens {
		tokenRepository.Add(&token)
	}

	service := GetUserTokensService{
		tokenRepository,
	}

	response, err := service.Execute(user.ID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Tokens) != 3 {
		t.Errorf("Expected to have %d tokens, but got %d", 3, len(response.Tokens))
	}
	if response.UserID != user.ID {
		t.Errorf("Expected user id to be %d, but got %d", user.ID, response.UserID)
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

	user := &u.User{
		Name:     name,
		Email:    email,
		Password: password,
		Consent:  consent,
	}
	id, err := userRepository.Add(user)
	if err != nil {
		t.Error(err)
	}
	user.ID = id

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

	response, err := service.Execute(user.ID)
	if err != nil {
		t.Error(err)
	}

	if len(response.Tokens) != 0 {
		t.Errorf("Expected to have %d tokens, but got %d", 0, len(response.Tokens))
	}
	if response.UserID != user.ID {
		t.Errorf("Expected user id to be %d, but got %d", user.ID, response.UserID)
	}

	t.Cleanup(func() {
		user := u.User{ID: response.UserID}
		userRepository.Remove(&user)

		tokenRepository.RemoveByUser(response.UserID)
	})
}
