package services

import (
	"testing"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	u "github.com/MihaiBlebea/Wordpress/platform/user"
)

func TestUserCannotLoginWhenNotRegistered(t *testing.T) {
	var (
		email    = "serban@gmail.com"
		password = "intrex"
	)

	userRepository := *u.Repo(&c.MysqlConnection{
		Username: "admin",
		Password: "pass",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "platform",
	})
	service := LoginUserService{
		userRepository,
	}
	response, err := service.Execute(email, password)
	if err != nil {
		t.Error(err)
	}
	if response.Success != false {
		t.Errorf("Response Success should be false")
	}
}

func TestUserTryToLogin(t *testing.T) {
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

	loginService := LoginUserService{
		userRepository,
	}

	t.Run("WithWrongEmail", func(t *testing.T) {
		response, err := loginService.Execute(email+"s", password)
		if err != nil {
			t.Error(err)
		}
		if response.Success != false {
			t.Errorf("Response Success should be false")
		}
		if len(response.JWT) != 0 {
			t.Errorf("JWT token length should be %d", 0)
		}
	})

	t.Run("WithWrongPassword", func(t *testing.T) {
		response, err := loginService.Execute(email, password+"s")
		if err != nil {
			t.Error(err)
		}
		if response.Success != false {
			t.Errorf("Response Success should be false")
		}
		if len(response.JWT) != 0 {
			t.Errorf("JWT token length should be %d", 0)
		}
	})

	t.Run("WithCorrectCredentials", func(t *testing.T) {
		response, err := loginService.Execute(email, password)
		if err != nil {
			t.Error(err)
		}
		if response.Success != true {
			t.Errorf("Response Success should be false")
		}
		if response.Message != "Authentication successfull" {
			t.Errorf("Response Message should be %s, but got %s", "Authentication successfull", response.Message)
		}
		if response.ID == 0 {
			t.Errorf("Response ID should not be 0")
		}
		if len(response.JWT) != 127 {
			t.Errorf("JWT token length should be %d", 127)
		}
	})

	t.Cleanup(func() {
		user := u.User{ID: registerResponse.ID}
		userRepository.Remove(&user)
	})
}
