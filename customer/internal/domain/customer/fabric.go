package customer

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
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

func CreateCustomer(name, phone, password string) (*Customer, error) {
	if !validateName(name) {
		return nil, ErrInvalidCustomerName
	}
	if !validatePhone(phone) {
		return nil, ErrInvalidCustomerPhone
	}
	if !validatePassword(password) {
		return nil, ErrInvalidCustomerPassword
	}

	customer := &Customer{
		ID:      uuid.New(),
		Name:    name,
		Phone:   phone,
		Created: time.Now(),
	}

	if err := customer.SetPassword(password); err != nil {
		return nil, err
	}

	return customer, nil
}
