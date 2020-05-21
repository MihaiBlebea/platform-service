package user

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	JWT      string    `json:"jwt"`
	Active   bool      `json:"active"`
	Consent  bool      `json:"consent"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Username returns a combination of firstname and lastname
func (u *User) Username() string {
	return fmt.Sprintf("%s", u.Name)
}

// HashPassword encrypts the password
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// GenerateJWT generates a JWT token on the user model
func (u *User) GenerateJWT() error {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = u.Name
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	u.JWT = signedToken
	return nil
}

// New returns new User model
func New(name, email, password string, consent bool) (*User, error) {
	user := &User{
		Name:    name,
		Email:   email,
		Active:  true,
		Consent: consent,
		Created: time.Now(),
		Updated: time.Now(),
	}
	err := user.HashPassword(password)
	if err != nil {
		return &User{}, err
	}
	err = user.GenerateJWT()
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// NewFromRaw returns a User model from database
func NewFromRaw(id int, name, email, password string, active, consent bool, created, updated time.Time) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Password: password,
		Active:   active,
		Consent:  consent,
		Created:  created,
		Updated:  updated,
	}
}
