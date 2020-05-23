package token

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const timeFormat = "2006-01-02 15:04:05"

// Token model
// https://stackoverflow.com/questions/45267125/how-to-generate-unique-random-alphanumeric-tokens-in-golang
type Token struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Token      string    `json:"token"`
	Hash       string    `json:"hash"`
	Active     bool      `json:"active"`
	LinkedHost string    `json:"linked_host"`
	Expires    string    `json:"expires"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

// Validate if the token is valid or not return bool
func (t *Token) Validate(token, host string) (bool, error) {
	if t.Active != true {
		return false, errors.New("Token is not active")
	}

	if t.Token != token {
		return false, errors.New("Token is invalid")
	}

	if t.LinkedHost != host {
		return false, errors.New("Token's host is not the same as the provided one")
	}

	expires, err := time.Parse(timeFormat, t.Expires)
	if err != nil {
		return false, err
	}
	now := time.Now()
	if expires.Before(now) {
		return false, errors.New("Token has expired")
	}
	return true, nil
}

// HasHost returns bool if token has linked host or not
func (t *Token) HasHost() bool {
	if t.LinkedHost == "" {
		return false
	}
	return true
}

// Link links the token with a host
func (t *Token) Link(host string) {
	t.LinkedHost = host
}

// Unlink detaches the host from the token
func (t *Token) Unlink() {
	t.LinkedHost = ""
}

// GenerateToken returns a string token
func GenerateToken(email string) (string, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil)), string(hash), nil
}

// New generate a new token for a user
func New(userID int, email string) (*Token, error) {
	expires := time.Now().Add(24 * 30 * time.Hour)

	token, hash, err := GenerateToken(email)
	if err != nil {
		return &Token{}, err
	}
	return &Token{
		UserID:     userID,
		Token:      token,
		Hash:       hash,
		Active:     true,
		LinkedHost: "",
		Expires:    expires.Format(timeFormat),
		Created:    time.Now(),
		Updated:    time.Now(),
	}, nil
}

// NewFromRaw returns a Token from raw data
func NewFromRaw(id, userID int, token, hash string, active bool, host string, expires string, created, updated time.Time) (*Token, error) {
	return &Token{
		ID:         id,
		UserID:     userID,
		Token:      token,
		Hash:       hash,
		Active:     active,
		LinkedHost: host,
		Expires:    expires,
		Created:    created,
		Updated:    updated,
	}, nil
}

// NewCount returns a number of tokens
func NewCount(userID int, email string, count int) (*[]Token, error) {
	tokens := []Token{}
	for i := 0; i < count; i++ {
		token, err := New(userID, email)
		if err != nil {
			return &tokens, err
		}
		tokens = append(tokens, *token)
	}
	return &tokens, nil
}
