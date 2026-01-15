package utils

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.New("Email is required")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("Invalid email format")
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("Password is required")
	}

	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters")
	}

	if len(password) > 72 {
		return errors.New("Password must be at most 72 characters")
	}

	return nil
}
