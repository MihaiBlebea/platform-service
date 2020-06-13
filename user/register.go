package user

import (
	c "github.com/MihaiBlebea/purpletree/platform/connection"
)

// Register registers a user and generates a JWT token
func Register(name, email, password string, consent bool) (user *User, err error) {
	userRepo := Repo(c.Mysql())

	user, err = New(name, email, password, consent)
	if err != nil {
		return user, err
	}

	userID, err := userRepo.Add(user)
	if err != nil {
		return user, err
	}
	user.ID = userID

	return user, nil
}
