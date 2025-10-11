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
	Email    string
	Password []byte
	Created  time.Time

	FailedCount        int
	LockedUntil        *time.Time
	PasswordUpdated    time.Time
	MustChangePassword bool
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
	c.PasswordUpdated = time.Now()

	return nil
}

func (c *Customer) SetNewPassword(newPassword string) error {
	if err := c.SetPassword(newPassword); err != nil {
		return err
	}
	c.MustChangePassword = false
	return nil
}

func (c *Customer) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(c.Password, []byte(password)) == nil
}

func (c *Customer) IsLocked() bool {
	return c.LockedUntil != nil && time.Now().Before(*c.LockedUntil)
}

func (c *Customer) RegisterFailedAttempt(lp LockoutPolicy) {
	c.FailedCount++
	if c.FailedCount >= lp.MaxFailed {
		until := time.Now().Add(lp.LockFor)
		c.LockedUntil = &until
	}
}

func (c *Customer) ResetFailedAttempts() {
	c.FailedCount = 0
	c.LockedUntil = nil
}

func (c *Customer) MarkMustChangePassword() {
	c.MustChangePassword = true
}

type LockoutPolicy struct {
	MaxFailed int
	LockFor   time.Duration
}
