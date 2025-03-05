package domain

import (
	"github.com/google/uuid"
	"regexp"
	"strings"
	"time"
)

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)

func validatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func validateName(name string) bool {
	return strings.TrimSpace(name) != ""
}

func validatePassword(password string) bool {
	return password != ""
}

func Create(name, phone, password string) (*Courier, error) {
	if !validateName(name) {
		return nil, ErrInvalidCourierName
	}
	if !validatePhone(phone) {
		return nil, ErrInvalidCourierPhone
	}
	if !validatePassword(password) {
		return nil, ErrInvalidCourierPassword
	}

	c := &Courier{
		ID:      uuid.New(),
		Name:    name,
		Phone:   phone,
		Created: time.Now(),
	}
	if err := c.SetPassword(password); err != nil {
		return nil, err
	}
	return c, nil
}
