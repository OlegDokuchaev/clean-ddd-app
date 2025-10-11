package customer

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

var phoneRegex = regexp.MustCompile(`^\+?[0-9]{10,15}$`)
var emailRegex = regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)

func validatePhone(phone string) bool {
	return phoneRegex.MatchString(phone)
}

func validateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func validateName(name string) bool {
	return strings.TrimSpace(name) != ""
}

func validatePassword(password string) bool {
	return password != ""
}

func Create(name, phone, email, password string) (*Customer, error) {
	if !validateName(name) {
		return nil, ErrInvalidCustomerName
	}
	if !validatePhone(phone) {
		return nil, ErrInvalidCustomerPhone
	}
	if !validateEmail(email) {
		return nil, ErrInvalidCustomerEmail
	}
	if !validatePassword(password) {
		return nil, ErrInvalidCustomerPassword
	}

	now := time.Now()

	customer := &Customer{
		ID:          uuid.New(),
		Name:        name,
		Phone:       phone,
		Email:       email,
		Created:     now,
		FailedCount: 0,
		LockedUntil: nil,
	}

	if err := customer.SetPassword(password); err != nil {
		return nil, err
	}

	return customer, nil
}
