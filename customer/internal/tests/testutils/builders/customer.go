package builders

import (
	"time"

	customerDomain "customer/internal/domain/customer"
	"github.com/google/uuid"
)

type CustomerBuilder struct {
	id    uuid.UUID
	name  string
	phone string
	email string

	passwordPlain  string
	passwordHashed []byte

	created     time.Time
	failedCount int
	lockedUntil *time.Time
}

func NewCustomerBuilder() *CustomerBuilder {
	return &CustomerBuilder{
		id:            uuid.New(),
		name:          "John Doe",
		phone:         "+79991234567",
		email:         "john.doe@example.com",
		passwordPlain: "P@ssw0rd!",
		created:       time.Now(),
		failedCount:   0,
		lockedUntil:   nil,
	}
}

func (b *CustomerBuilder) WithID(id uuid.UUID) *CustomerBuilder {
	b.id = id
	return b
}

func (b *CustomerBuilder) WithName(name string) *CustomerBuilder {
	b.name = name
	return b
}

func (b *CustomerBuilder) WithPhone(phone string) *CustomerBuilder {
	b.phone = phone
	return b
}

func (b *CustomerBuilder) WithEmail(email string) *CustomerBuilder {
	b.email = email
	return b
}

func (b *CustomerBuilder) WithPassword(plain string) *CustomerBuilder {
	b.passwordPlain = plain
	return b
}

func (b *CustomerBuilder) WithHashedPassword(hashed []byte) *CustomerBuilder {
	b.passwordHashed = hashed
	return b
}

func (b *CustomerBuilder) WithCreated(t time.Time) *CustomerBuilder {
	b.created = t
	return b
}

func (b *CustomerBuilder) WithFailedCount(n int) *CustomerBuilder {
	b.failedCount = n
	return b
}

func (b *CustomerBuilder) LockedUntil(t time.Time) *CustomerBuilder {
	b.lockedUntil = &t
	return b
}

func (b *CustomerBuilder) LockedFor(d time.Duration) *CustomerBuilder {
	t := time.Now().Add(d)
	b.lockedUntil = &t
	return b
}

func (b *CustomerBuilder) Unlocked() *CustomerBuilder {
	b.lockedUntil = nil
	b.failedCount = 0
	return b
}

func (b *CustomerBuilder) Build() *customerDomain.Customer {
	c := &customerDomain.Customer{
		ID:          b.id,
		Name:        b.name,
		Phone:       b.phone,
		Email:       b.email,
		Created:     b.created,
		FailedCount: b.failedCount,
		LockedUntil: b.lockedUntil,
	}

	if len(b.passwordHashed) > 0 {
		c.Password = make([]byte, len(b.passwordHashed))
		copy(c.Password, b.passwordHashed)
	} else {
		_ = c.SetPassword(b.passwordPlain)
	}

	return c
}
