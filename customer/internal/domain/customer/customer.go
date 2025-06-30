package customer

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	ID       uuid.UUID
	Name     string
	Phone    string
	Password []byte
	Created  time.Time
}

func (c *Customer) SetPassword(password string) error {
	if !validatePassword(password) {
		return ErrInvalidCustomerPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrInvalidCustomerPassword
	}
	c.Password = hash
	return nil
}

func (c *Customer) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(c.Password, []byte(password)) == nil
}
