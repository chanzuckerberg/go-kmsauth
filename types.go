package kmsauth

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	validator "gopkg.in/go-playground/validator.v9"
)

// ------------- AuthContext --------------

// AuthContext is a kms encryption context used to ascertain a user's identiy
type AuthContext interface {
	Validate() error
	GetUsername() string
	GetKMSContext() map[string]*string
}

// AuthContextV1 is a kms encryption context used to ascertain a user's identiy
type AuthContextV1 struct {
	From string `json:"from,omitempty" validate:"required"`
	To   string `json:"to,omitempty" validate:"required"`
}

// Validate validates
func (ac *AuthContextV1) Validate() error {
	if ac == nil {
		return errors.New("NilAuthContext")
	}
	v := validator.New()
	return v.Struct(ac)
}

// GetUsername returns a username
func (ac *AuthContextV1) GetUsername() string {
	return ac.From
}

// GetKMSContext gets the kms context
func (ac *AuthContextV1) GetKMSContext() map[string]*string {
	return map[string]*string{
		"from": &ac.From,
		"to":   &ac.To,
	}
}

// AuthContextV2 is a kms encryption context used to ascertain a user's identiy
type AuthContextV2 struct {
	From     string `json:"from,omitempty" validate:"required"`
	To       string `json:"to,omitempty" validate:"required"`
	UserType string `json:"user_type,omitempty" validate:"required"`
}

// Validate validates
func (ac *AuthContextV2) Validate() error {
	if ac == nil {
		return errors.New("NilAuthContext")
	}
	v := validator.New()
	return v.Struct(ac)
}

// GetUsername returns a username
func (ac *AuthContextV2) GetUsername() string {
	return fmt.Sprintf("%d/%s/%s", TokenVersion2, ac.UserType, ac.From)
}

// GetKMSContext gets the kms context
func (ac *AuthContextV2) GetKMSContext() map[string]*string {
	return map[string]*string{
		"from": &ac.From,
		"to":   &ac.To,
		"user": &ac.UserType,
	}
}

// ------------- Token --------------

// TokenTime is a custom time formatter
type TokenTime time.Time

// MarshalJSON marshals into json
func (t TokenTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(t).Format(timeFormat)), nil
}

// UnmarshalJSON unmarshals
func (t TokenTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = strings.Trim(s, "\"")
	parsed, err := time.Parse(timeFormat, s)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("Could not parse time %s", s))
	}
	t = TokenTime(parsed)
	return nil
}

// Token is a kmsauth token
type Token struct {
	NotBefore TokenTime `json:"not_before,omitempty"`
	NotAfter  TokenTime `json:"not_after,omitempty"`
}

// ------------- EncryptedToken --------------

// EncryptedToken is a b64 kms encrypted token
type EncryptedToken string
