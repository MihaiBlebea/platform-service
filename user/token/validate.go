package token

import (
	c "github.com/MihaiBlebea/Wordpress/platform/connection"
)

// Validate token
func Validate(token, host string) (isValid bool, err error) {
	tokenRepo := Repo(c.Mysql())
	tkn, _, err := tokenRepo.FindToken(token)
	if err != nil {
		return false, err
	}

	valid, err := tkn.Validate(token, host)
	if err != nil {
		return false, err
	}

	if valid != true {
		return false, err
	}

	return true, nil
}
