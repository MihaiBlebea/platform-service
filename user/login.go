package user

import (
	"errors"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
)

// Authenticate checks if a JWT is valid or not and returns bool
func Authenticate(token string) (bool, *User, error) {
	userRepo := Repo(c.Mysql())
	user, count, err := userRepo.FindByJWT(token)
	if err != nil {
		return false, &User{}, nil
	}
	if count == 0 {
		return false, &User{}, errors.New("No user found")
	}

	return true, user, nil
}
