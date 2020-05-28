package user

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	JWT         string    `json:"jwt"`
	Active      bool      `json:"active"`
	Consent     bool      `json:"consent"`
	ConfirmCode string    `json:"confirm_code"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
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

// GenerateRandomPassword generates a random password for the user
func GenerateRandomPassword() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
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
