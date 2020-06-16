package usrauth

import (
	"errors"

	c "github.com/MihaiBlebea/purpletree/platform/connection"
	u "github.com/MihaiBlebea/purpletree/platform/user"
)

// New returns a new Service
func New() *Service {
	userRepository := *u.Repo(c.Mysql())
	return &Service{userRepository}
}

// Service for authenticating user by jwt token
type Service struct {
	userRepository u.Repository
}

// Execute runs the Service
func (s *Service) Execute(jwt string) (isAuth bool, user *u.User, err error) {
	user, count, err := s.userRepository.FindByJWT(jwt)
	if err != nil {
		return false, user, err
	}
	if count == 0 {
		return false, user, errors.New("No user found in db by jwt token")
	}

	return true, user, nil
}
