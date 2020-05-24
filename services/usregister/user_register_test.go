package usregister

import (
	"testing"

	c "github.com/MihaiBlebea/Wordpress/platform/connection"
	u "github.com/MihaiBlebea/Wordpress/platform/user"
)

func TestUserCanRegister(t *testing.T) {
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
	service := RegisterUserService{
		userRepository,
	}
	response, err := service.Execute(name, email, password, consent)
	if err != nil {
		t.Error(err)
	}
	if response.Name != name {
		t.Errorf("Name %s should be equal to %s", response.Name, name)
	}
	if response.JWT == "" {
		t.Errorf("Response should have a JWT token")
	}
	if len(response.JWT) != 127 {
		t.Errorf("JWT token length should be %d", 127)
	}

	t.Cleanup(func() {
		user := u.User{ID: response.ID}
		userRepository.Remove(&user)
	})
}
