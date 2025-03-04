package domain

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Courier struct {
	ID       uuid.UUID
	Name     string
	Phone    string
	Password []byte
	Created  time.Time
}

func (c *Courier) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ErrInvalidCourierPassword
	}

	c.Password = hash
	return nil
}

func (c *Courier) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(c.Password, []byte(password))
	return err == nil
}
